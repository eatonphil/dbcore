package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"{{api.repo}}/go/pkg/dao"
)

type Server struct {
	dao *dao.DAO
	router *httprouter.Router
	logger logrus.FieldLogger
	address string
}

func (s Server) registerControllers() {
	{{~ for table in tables ~}}
	// Register {{table.name}} routes
	s.router.GET("/{{table.name}}", s.{{table.name}}GetManyController)
	s.router.POST("/{{table.name}}/new", s.{{table.name}}CreateController)
	{{~ if table.primary_key.is_some ~}}
	s.router.GET("/{{table.name}}/:key", s.{{table.name}}GetController)
	s.router.PUT("/{{table.name}}/:key", s.{{table.name}}UpdateController)
	s.router.DELETE("/{{table.name}}/:key", s.{{table.name}}DeleteController)
	{{~ end ~}}
	{{ end }}
}

func (s Server) registerSigintHandler(srv *http.Server) {
	// Wait for SIGINT
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	s.logger.Info("Signal received, shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Warnf("Error during shutdown: %s", err)
		return
	}
}

func (s Server) Start() {
	s.registerControllers()

	srv := &http.Server{
		Handler: loggerRouter{s.logger, s.router},
		Addr:    s.address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		s.logger.Infof("Starting server at %s", s.address)
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Warnf("Exiting server with error: %s", err)
			return
		}
		
		s.logger.Info("Exiting server")
	}()

	s.registerSigintHandler(srv)
}

func New(conf *Config) (*Server, error) {
	dialect := conf.GetString("database.dialect")
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s%s",
		dialect,
		conf.GetString("database.username"),
		conf.GetString("database.password"),
		conf.GetString("database.host", "localhost"),
		conf.GetString("database.port", "5432"),
		conf.GetString("database.database"),
		conf.GetString("database.parameters"))
	db, err := sqlx.Connect(dialect, dsn)
	if err != nil {
		return nil, err
	}

	router := httprouter.New()
	logger := logrus.New()

	return &Server{
		dao: dao.New(db),
		router: router,
		logger: logger.WithFields(logrus.Fields{
			"struct": "Server",
			"pkg": "server",
		}),
		address: conf.GetString("api.address", ":9090"),
	}, nil
}

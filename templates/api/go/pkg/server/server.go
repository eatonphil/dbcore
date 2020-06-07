package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	{{ if database.dialect == "postgres" }}
	_ "github.com/lib/pq"
	{{ else if database.dialect == "mysql" }}
	_ "github.com/go-sql-driver/mysql"
	{{ end }}
	"github.com/sirupsen/logrus"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
)

type Server struct {
	dao *dao.DAO
	router *httprouter.Router
	logger logrus.FieldLogger
	address string
	secret string
	sessionDuration time.Duration
}

func (s Server) registerControllers() {
	{{~ for table in tables ~}}
	// Register {{table.name}} routes
	s.router.GET("/{{ api.router_prefix }}{{ table.name }}", s.{{table.name}}GetManyController)
	s.router.POST("/{{ api.router_prefix }}{{ table.name }}/new", s.{{table.name}}CreateController)
	{{~ if table.primary_key.is_some ~}}
	s.router.GET("/{{ api.router_prefix }}{{ table.name }}/:key", s.{{table.name}}GetController)
	s.router.PUT("/{{ api.router_prefix }}{{ table.name }}/:key", s.{{table.name}}UpdateController)
	s.router.DELETE("/{{ api.router_prefix }}{{ table.name }}/:key", s.{{table.name}}DeleteController)
	{{~ end ~}}
	{{ end }}

	{{ if api.auth.enabled }}
	// Register session route
	s.router.POST("/{{ api.router_prefix }}session/start", s.SessionStartController)
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
		Handler: s,
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
	dsn := conf.GetString("database.dsn")
	db, err := sqlx.Connect(dialect, dsn)
	if err != nil {
		return nil, err
	}

	secret := conf.GetString("session.secret")
	{{ if api.auth.enabled }}
	if secret == "" {
		return nil, fmt.Errorf(`Configuration value "secret" must be specified`)
	}
	{{ end }}

	router := httprouter.New()
	logger := logrus.New()
	// TODO: make this configurable
	logger.SetLevel(logrus.DebugLevel)

	return &Server{
		dao: dao.New(db),
		router: router,
		logger: logger.WithFields(logrus.Fields{
			"struct": "Server",
			"pkg": "server",
		}),
		address: conf.GetString("address", ":9090"),
		secret: secret,
		sessionDuration: conf.GetDuration("session.duration", time.Hour * 2),
	}, nil
}

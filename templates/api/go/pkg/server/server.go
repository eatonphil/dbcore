package server

import (
	"context"
	"fmt"
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
	{{ else if database.dialect == "sqlite" }}
	_ "github.com/mattn/go-sqlite3"
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
	allowedOrigins []string
}

func (s Server) registerControllers() {
	{{~ for table in tables ~}}
	// Register {{table.label}} routes
	s.router.GET("/{{ api.router_prefix }}{{ table.label }}", s.{{table.label}}GetManyController)
	s.router.POST("/{{ api.router_prefix }}{{ table.label }}", s.{{table.label}}CreateController)
	{{~ if table.primary_key.value ~}}
	s.router.GET("/{{ api.router_prefix }}{{ table.label }}/:key", s.{{table.label}}GetController)
	s.router.PUT("/{{ api.router_prefix }}{{ table.label }}/:key", s.{{table.label}}UpdateController)
	s.router.DELETE("/{{ api.router_prefix }}{{ table.label }}/:key", s.{{table.label}}DeleteController)
	{{~ end ~}}
	{{ end }}

	{{ if api.auth.enabled }}
	// Register session route
	s.router.POST("/{{ api.router_prefix }}session/start", s.SessionStartController)
	s.router.POST("/{{ api.router_prefix }}session/stop", s.SessionStopController)
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

func (s Server) handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
	s.logger.Warnf("Unexpected panic: %s", err)
	sendErrorResponse(w, fmt.Errorf("Internal server error"))
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Infof("%s %s", r.Method, r.URL.RequestURI())
	w.Header().Set("Content-Type", "application/json")

	for _, allowed := range s.allowedOrigins {
		origin := strings.ToLower(r.Header.Get("origin"))
		if strings.ToLower(allowed) == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
		}

		if r.Method == http.MethodOptions {
			return
		}
	}

	s.router.ServeHTTP(w, r)
}


func (s Server) Start() {
	s.router.PanicHandler = s.handlePanic
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
		dao: dao.New(db, logger),
		router: router,
		logger: logger.WithFields(logrus.Fields{
			"struct": "Server",
			"pkg": "server",
		}),
		address: conf.GetString("address", ":9090"),
		secret: secret,
		sessionDuration: conf.GetDuration("session.duration", time.Hour * 2),
		allowedOrigins : conf.GetStringSlice("allowed-origins"),
	}, nil
}

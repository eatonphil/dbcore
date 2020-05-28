package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type loggerRouter struct {
	logger logrus.FieldLogger
	srv http.Handler
}

func (lr loggerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lr.logger.Infof("%s %s", r.Method, r.URL.RequestURI())
	lr.srv.ServeHTTP(w, r)
}

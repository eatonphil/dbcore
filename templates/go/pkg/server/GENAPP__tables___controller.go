package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s Server) {{table}}GetManyController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.dao.{{table | string.capitalize}}GetMany()
}

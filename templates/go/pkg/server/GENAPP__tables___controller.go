package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"{{api.repo}}/pkg/dao"
)

func (s Server) {{table.name}}GetManyController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result, err := s.dao.{{table.name|string.capitalize}}GetMany()
	if err != nil {
		// TODO: handle error
		sendErrorResponse(w, err)
		return
	}

	sendPaginatedResponse(w, result)
}

func (s Server) {{table.name}}GetController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getIntParameter(ps, "id")
	if err != nil {
		sendValidationErrorsResponse(map[string]string[]{
			"id", []string{"Expected integer id"},
		}, err)
		return
	}

	result, err := s.dao.{{table.name|string.capitalize}}Get(id)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, result)
}

func (s Server) {{table.name}}UpdateController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getIntParameter(ps, "id")
	if err != nil {
		sendValidationErrorsResponse(map[string]string[]{
			"id": []string{"Expected integer id"},
		})
		return
	}

	var body dao.{{table.name|string.capitalize}}
	err = getBody(&body)
	if err != nil {
		sendValidationErrorResponse(w, "Expected valid JSON", err)
		return
	}

	result, err := s.dao.{{table.name|string.capitalize}}Update(id, body)
	if err != nil {
		// TODO: handle error
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, result)
}

func (s Server) {{table.name}}CreateController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body dao.{{table.name|string.capitalize}}
	err = getBody(&body)
	if err != nil {
		sendValidationErrorResponse(w, "Expected valid JSON", err)
		return
	}

	result, err := s.dao.{{table.name|string.capitalize}}Insert(body)
	if err != nil {
		// TODO: handle error
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, result)
}

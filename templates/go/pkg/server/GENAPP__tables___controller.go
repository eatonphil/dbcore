package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"{{ api.repo }}/pkg/dao"
)

func (s Server) {{ table.name }}GetManyController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter, pageInfo, err := getFilterAndPageInfo(r)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	result, err := s.dao.{{ table.name|string.capitalize }}GetMany(filter, *pageInfo)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, result)
}

func (s Server) {{ table.name }}CreateController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body dao.{{ table.name|string.capitalize }}
	err := getBody(r, &body)
	if err != nil {
		s.logger.Debug("Expected valid JSON, got: %s", err)
		sendValidationErrorResponse(w, "Expected valid JSON")
		return
	}

	err := s.dao.{{ table.name|string.capitalize }}Insert(&body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, body)
}

{{ if table.primary_key }}
func (s Server) {{ table.name }}GetController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result, err := s.dao.{{ table.name|string.capitalize }}Get(ps.ByName("key"))
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, result)
}

func (s Server) {{ table.name }}UpdateController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var body dao.{{ table.name|string.capitalize }}
	err = getBody(r, &body)
	if err != nil {
		sendValidationErrorResponse(w, "Expected valid JSON", err)
		return
	}

	result, err := s.dao.{{ table.name|string.capitalize }}Update(ps.ByName("key"), body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, result)
}
{{ end }}

package server

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
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

	{{ if table.name == api.auth.table }}
	for i, _ := range result.Data {
		// TODO: make sure this column actually exists
		result.Data[i].C_{{ api.auth.password }} = "<REDACTED>"
	}
	{{ end }}

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

	{{ if api.auth.enabled && table.name == api.auth.table }}
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(body.C_{{ api.auth.password }}), bcrypt.DefaultCost)
	body.C_{{ api.auth.password }} = string(hash)
	{{ end }}

	err = s.dao.{{ table.name|string.capitalize }}Insert(&body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.name == api.auth.table }}
	// TODO: make sure this column actually exists
	body.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, body)
}

{{ if table.primary_key.enabled }}
func (s Server) {{ table.name }}GetController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	result, err := s.dao.{{ table.name|string.capitalize }}Get(ps.ByName("key"))
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.name == api.auth.table }}
	// TODO: make sure this column actually exists
	result.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, result)
}

func (s Server) {{ table.name }}UpdateController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var body dao.{{ table.name|string.capitalize }}
	err := getBody(r, &body)
	if err != nil {
		s.logger.Debug("Expected valid JSON, got: %s", err)
		sendValidationErrorResponse(w, "Expected valid JSON")
		return
	}

	result, err := s.dao.{{ table.name|string.capitalize }}Get(ps.ByName("key"))
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if api.auth.enabled && table.name == api.auth.table }}
	body.C_{{ api.auth.password }} = result.C_{{ api.auth.password }}
	{{ end }}

	result, err = s.dao.{{ table.name|string.capitalize }}Update(ps.ByName("key"), body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.name == api.auth.table }}
	// TODO: make sure this column actually exists
	result.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, result)
}
{{ end }}

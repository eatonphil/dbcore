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
	body.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, body)
}

{{~ if table.primary_key.value ~}}
{{~
  func toGoType
    case $0
      when "int"
        "int"
      when "bigint"
        "int64"
      when "text", "varchar", "char"
        "string"
      when "boolean"
        "bool"
      when "timestamp"
        "time.Time"
      else
        "Unsupported PostgreSQL type: " + $0.type
    end
  end
~}}

func parse{{ table.name|string.capitalize }}Key(key string) {{ toGoType table.primary_key.value.type }} {
{{~
  case table.primary_key.value.type
    when "text", "varchar", "char"
      "\t return key"
    when "int", "bigint"
      "\t i, _ := strconv.ParseInt(key, 10, 64)\n\t return i"
    when "timestamp"
      "\t t, _ := time.Parse(time.RFC3339, key)\n\t return t"
    when "boolean"
      "\t return key == \"true\""
    else
      "\t return \"\""
  end
~}}
}

func (s Server) {{ table.name }}GetController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	k := parse{{ table.name|string.capitalize }}Key(ps.ByName("key"))
	result, err := s.dao.{{ table.name|string.capitalize }}Get(k)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.name == api.auth.table }}
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

	k := parse{{ table.name|string.capitalize }}Key(ps.ByName("key"))
	{{ if api.auth.enabled && table.name == api.auth.table }}
	result, err := s.dao.{{ table.name|string.capitalize }}Get(k)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	body.C_{{ api.auth.password }} = result.C_{{ api.auth.password }}
	{{ end }}

	body.C_{{ table.primary_key.value.column }} = k
	err = s.dao.{{ table.name|string.capitalize }}Update(k, body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.name == api.auth.table }}
	body.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, body)
}

func (s Server) {{ table.name }}DeleteController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	k := parse{{ table.name|string.capitalize }}Key(ps.ByName("key"))
	err := s.dao.{{ table.name|string.capitalize }}Delete(k)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, struct{}{})
}
{{~ end ~}}

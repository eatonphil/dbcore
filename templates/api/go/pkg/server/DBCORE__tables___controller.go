package server

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
)

{{~ if table.primary_key.value ~}}
{{~
  func toGoType
    case $0.type
      when "int", "integer"
        "int32"
      when "bigint"
        "int64"
      when "text", "varchar", "char"
        "string"
      when "boolean"
        "bool"
      when "timestamp", "timestamp with time zone"
        "time.Time"
      else
        "Unsupported type: " + $0.type
    end
  end
~}}

func (s Server) {{ table.label }}RequestFilterContext(
	r *http.Request,
	objectId *{{ toGoType table.primary_key.value }},
) map[string]interface{} {
	ctx := map[string]interface{}{
		"req_username": s.getSessionUsername(r),
	}

	if objectId != nil {
		ctx["req_object_id"] = *objectId
	}

	return ctx
}

func (s Server) {{ table.label }}RequestIsAllowed(
	r *http.Request,
	filter string,
	objectId *{{ toGoType table.primary_key.value }},
) bool {
	ctx := s.{{ table.label }}RequestFilterContext(r, objectId)
	return s.dao.{{ table.label|dbcore_capitalize }}IsAllowed(filter, ctx)
}

func (s Server) {{ table.label }}GetManyController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	extraFilter, pageInfo, err := getFilterAndPageInfo(r)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	baseFilter := s.auth.allow["{{ table.label }}"]["get"]
	baseContext := s.{{ table.label }}RequestFilterContext(r, nil)

	result, err := s.dao.{{ table.label|dbcore_capitalize }}GetMany(extraFilter, *pageInfo, baseFilter, baseContext)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.label == api.auth.table }}
	for i, _ := range result.Data {
		result.Data[i].C_{{ api.auth.password }} = "<REDACTED>"
	}
	{{ end }}

	sendResponse(w, result)
}

func (s Server) {{ table.label }}CreateController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	baseFilter := s.auth.allow["{{ table.label }}"]["post"]
	if baseFilter != "" {
		if !s.{{ table.label }}RequestIsAllowed(r, baseFilter, nil) {
			sendAuthorizationErrorResponse(w)
			return
		}
	}

	var body dao.{{ table.label|dbcore_capitalize }}
	err := getBody(r, &body)
	if err != nil {
		s.logger.Debug("Expected valid JSON, got: %s", err)
		sendValidationErrorResponse(w, "Expected valid JSON")
		return
	}

	{{ if api.auth.enabled && table.label == api.auth.table }}
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(body.C_{{ api.auth.password }}), bcrypt.DefaultCost)
	body.C_{{ api.auth.password }} = string(hash)
	{{ end }}

	err = s.dao.{{ table.label|dbcore_capitalize }}Insert(&body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.label == api.auth.table }}
	body.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, body)
}

func parse{{ table.label|dbcore_capitalize }}Key(key string) {{ toGoType table.primary_key.value }} {
{{~
  case table.primary_key.value.type
    when "text", "varchar", "char"
      "\t return key"
    when "int", "integer"
      "\t i, _ := strconv.ParseInt(key, 10, 32)\n\t return int32(i)"
    when "bigint"
      "\t i, _ := strconv.ParseInt(key, 10, 64)\n\t return i"
    when "timestamp", "timestamp with time zone"
      "\t t, _ := time.Parse(time.RFC3339, key)\n\t return t"
    when "boolean"
      "\t return key == \"true\""
    else
      "\t return \"\""
  end
~}}
}

func (s Server) {{ table.label }}GetController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	k := parse{{ table.label|dbcore_capitalize }}Key(ps.ByName("key"))

	baseFilter := s.auth.allow["{{ table.label }}"]["get"]
	if baseFilter != "" {
		if !s.{{ table.label }}RequestIsAllowed(r, baseFilter, &k) {
			sendAuthorizationErrorResponse(w)
			return
		}
	}

	result, err := s.dao.{{ table.label|dbcore_capitalize }}Get(k)
	if err != nil {
		if err == dao.ErrNotFound {
			sendNotFoundErrorResponse(w)
			return
		}

		sendErrorResponse(w, err)
		return
	}

	{{ if table.label == api.auth.table }}
	result.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, result)
}

func (s Server) {{ table.label }}UpdateController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	k := parse{{ table.label|dbcore_capitalize }}Key(ps.ByName("key"))

	baseFilter := s.auth.allow["{{ table.label }}"]["put"]
	if baseFilter != "" {
		if !s.{{ table.label }}RequestIsAllowed(r, baseFilter, &k) {
			sendAuthorizationErrorResponse(w)
			return
		}
	}

	var body dao.{{ table.label|dbcore_capitalize }}
	err := getBody(r, &body)
	if err != nil {
		s.logger.Debug("Expected valid JSON, got: %s", err)
		sendValidationErrorResponse(w, "Expected valid JSON")
		return
	}

	{{ if api.auth.enabled && table.label == api.auth.table }}
	result, err := s.dao.{{ table.label|dbcore_capitalize }}Get(k)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	body.C_{{ api.auth.password }} = result.C_{{ api.auth.password }}
	{{ end }}

	body.C_{{ table.primary_key.value.column }} = k
	err = s.dao.{{ table.label|dbcore_capitalize }}Update(k, body)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	{{ if table.label == api.auth.table }}
	body.C_{{ api.auth.password }} = "<REDACTED>"
	{{ end }}

	sendResponse(w, body)
}

func (s Server) {{ table.label }}DeleteController(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	k := parse{{ table.label|dbcore_capitalize }}Key(ps.ByName("key"))

	baseFilter := s.auth.allow["{{ table.label }}"]["delete"]
	if baseFilter != "" {
		if !s.{{ table.label }}RequestIsAllowed(r, baseFilter, &k) {
			sendAuthorizationErrorResponse(w)
			return
		}
	}

	err := s.dao.{{ table.label|dbcore_capitalize }}Delete(k)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	sendResponse(w, struct{}{})
}
{{~ end ~}}

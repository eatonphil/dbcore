package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
)

func sendAuthorizationErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(struct{
		Error string `json:"error"`
	}{
		"Restricted interaction",
	})
}

func sendNotFoundErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(struct{
		Error string `json:"error"`
	}{
		"Not found",
	})}


func sendValidationErrorResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct{
		Error string `json:"error"`
	}{
		msg,
	})
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(struct{
		Error string `json:"error"`
	}{
		err.Error(),
	})
}

func sendResponse(w http.ResponseWriter, obj interface{}) {
	json.NewEncoder(w).Encode(obj)
}

func getBody(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(obj)
}

func getFilterAndPageInfo(r *http.Request) (*dao.Filter, *dao.Pagination, error) {
	getSingleUintParameter := func(param string, def uint64) (uint64, error) {
		values, ok := r.URL.Query()[param]
		if !ok || len(values) == 0 {
			return def, nil
		}

		return strconv.ParseUint(values[0], 10, 64)
	}

	limit, err := getSingleUintParameter("limit", 25)
	if err != nil {
		return nil, nil, err
	}

	offset, err := getSingleUintParameter("offset", 0)
	if err != nil {
		return nil, nil, err
	}

	sortColumn := r.URL.Query().Get("sortColumn")
	if sortColumn == "" {
		return nil, nil, fmt.Errorf(`Expected "sortColumn" parameter`)
	}

	sortOrder := strings.ToUpper(r.URL.Query().Get("sortOrder"))
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	if !(sortOrder == "ASC" || sortOrder == "DESC") {
		return nil, nil, fmt.Errorf(`Expected "sortOrder" parameter to be "asc" or "desc"`)
	}

	var filter *dao.Filter
	filterString := r.URL.Query().Get("filter")
	if filterString != "" {
		filter, err = dao.ParseFilter(filterString)
		if err != nil {
			return nil, nil, fmt.Errorf(`Expected valid "filter" parameter: %s`, err)
		}
	}

	return filter, &dao.Pagination{
		Limit: limit,
		Offset: offset,
		Order: sortColumn + " " + sortOrder,
	}, nil
}

func (s Server) getSessionUsername(r *http.Request) string {
	cookie, err := r.Cookie("au")
	if err != nil {
		// Fall back to header check
		cookie = &http.Cookie{}
	}

	token := cookie.Value
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > len("bearer ") &&
			strings.ToLower(authHeader[:len("bearer ")]) == "bearer " {
			token = authHeader[len("bearer "):]
		}
	}

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.secret), nil
	})
	if err != nil {
		s.logger.Debugf("Error parsing JWT: %s", err)
		return ""
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !(ok && t.Valid) {
		return ""
	}

	if err := claims.Valid(); err != nil {
		s.logger.Debugf("Invalid claims: %s", err)
		return ""
	}

	usernameInterface, ok := claims["username"]
	if !ok {
		return ""
	}

	username, ok := usernameInterface.(string)
	if !ok {
		return ""
	}

	return username
}


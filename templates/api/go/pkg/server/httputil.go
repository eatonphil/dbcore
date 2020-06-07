package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
)

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
	getSingleUintParameter := func(param string) (uint64, error) {
		values, ok := r.URL.Query()[param]
		if !ok || len(values) == 0 {
			return 0, fmt.Errorf(`Expected "%s" parameter`, param)
		}

		return strconv.ParseUint(values[0], 10, 64)
	}

	limit, err := getSingleUintParameter("limit")
	if err != nil {
		return nil, nil, err
	}

	offset, err := getSingleUintParameter("offset")
	if err != nil {
		return nil, nil, err
	}

	sortColumn := r.URL.Query().Get("sortColumn")
	if sortColumn == "" {
		return nil, nil, fmt.Errorf(`Expected "sortColumn" parameter`)
	}

	sortOrder := strings.ToLower(r.URL.Query().Get("sortOrder"))
	if !(sortOrder == "asc" || sortOrder == "desc") {
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

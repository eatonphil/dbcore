package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"

	"{{ api.repo }}/pkg/dao"
)

func sendValidationErrorResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(struct{
		Error string `json:"error"`
	}{
		msg,
	})
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(struct{
		Error string `json:"error"`
	}{
		err.Error(),
	})
}

func sendResponse(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func getBody(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(obj)
}

func getFilterAndPageInfo(r *http.Request) (squirrel.Sqlizer, *dao.Pagination, error) {
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

	// TODO: support actual squirrel filters
	return nil, &dao.Pagination{
		Limit: limit,
		Offset: offset,
		Order: sortColumn + " " + sortOrder,
	}, nil
}

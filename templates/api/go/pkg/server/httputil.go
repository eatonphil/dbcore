package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/xwb1989/sqlparser"
	"github.com/xwb1989/sqlparser/dependency/querypb"

	"{{ api.extra.repo }}/go/pkg/dao"
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

func getFilterAndPageInfo(r *http.Request) (*sqlparser.Expr, *dao.Pagination, error) {
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

	var parsedFilter *sqlparser.Expr
	filter := r.URL.Query().Get("filter")
	if filter != "" {
		stmt, err := sqlparser.Parse("SELECT 1 WHERE " + filter)
		if err != nil {
			return nil, nil, fmt.Errorf(`Expected valid "filter" parameter: %s`, err)
		}

		bv := map[string]*querypb.BindVariable{}
		sqlparser.Normalize(stmt, bv, "")
		fmt.Println(bv)
		//parsedFilter = stmt.Where.Expr
	}

	return parsedFilter, &dao.Pagination{
		Limit: limit,
		Offset: offset,
		Order: sortColumn + " " + sortOrder,
	}, nil
}

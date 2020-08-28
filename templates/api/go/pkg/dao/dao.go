package dao

import (
	"errors"
	"regexp"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/xwb1989/sqlparser"
	"github.com/xwb1989/sqlparser/dependency/querypb"
	"github.com/xwb1989/sqlparser/dependency/sqltypes"
)

var ErrNotFound = errors.New("Not found")

type Pagination struct {
	Limit uint64
	Offset uint64
	Order string
}


type Filter struct {
	args []interface{}
	filter string
}

func ParseFilter(filter string) (*Filter, error) {
	if filter == "" {
		return nil, nil
	}

	// TODO: validate filter uses acceptable subset of WHERE

	// Add stub select to make filter into a statement
	stmt, err := sqlparser.Parse("SELECT 1 WHERE " + filter)
	if err != nil {
		return nil, err
	}

	where := stmt.(*sqlparser.Select).Where

	bindings := map[string]*querypb.BindVariable{}
	// dbcore is a random choice for a binding variable prefix so
	// we can be more certain when we find-replace on it
	sqlparser.Normalize(stmt, bindings, "dbcore")

	// Map :dbcore[0-9]* to $[0-9]*, will screw up if this literal
	// appeared in text in the query.
	re := regexp.MustCompile(":dbcore([0-9]*)")

	var invalidValue error
	var args []interface{}
	whereWithBindings := re.ReplaceAllStringFunc(sqlparser.String(where), func (match string) string {
		// This library has no sane way to produce a Go value
		// from a parsed bind variable.
		match = match[1:] // Drop the preceeding colon
		v, _ := sqltypes.BindVariableToValue(bindings[match])
		s := v.ToString()
		if v.IsSigned() {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				invalidValue = err
				return ""
			}

			args = append(args, i)
		} else if v.IsFloat() {
			fl, err := strconv.ParseFloat(s, 64)
			if err != nil {
				invalidValue = err
				return ""
			}

			args = append(args, fl)
		} else if v.IsText() || v.IsQuoted() {
			args = append(args, s)
		} else if v.IsNull() {
			args = append(args, nil)
		} else {
			invalidValue = fmt.Errorf(`Unsupported value: "%s"`, s)
		}

		{{ if database.dialect == "postgres" }}
		return fmt.Sprintf("$%d", len(args))
		{{ else if database.dialect == "mysql" || database.dialect == "sqlite" }}
		return "?"
		{{ end }}
	})

	return &Filter{
		// Take only the filter part from the statement
		filter: whereWithBindings,
		args: args,
	}, nil
}

// map $-prefixed variables from request context that can be turned
// into parameterized queries
func applyVariablesFromContext(filter string, ctx map[string]interface{}) (string, []interface{}) {
	re := regexp.MustCompile("\\$[a-zA-Z_]+")
	var args []interface{}
	applied := re.ReplaceAllStringFunc(filter, func (match string) string {
		mapping, ok := ctx[match[1:]] // skip $ prefix
		if mapping == nil || !ok {
			args = append(args, nil)
		} else {
			args = append(args, mapping)			
		}

		{{ if database.dialect == "postgres" }}
		return fmt.Sprintf("$%d", len(args))
		{{ else if database.dialect == "mysql" || database.dialect == "sqlite" }}
		return "?"
		{{ end }}
	})

	return applied, args
}

type DAO struct {
	db *sqlx.DB
	logger logrus.FieldLogger
}

func New(db *sqlx.DB, logger logrus.FieldLogger) *DAO {
	return &DAO{
		db: db,
		logger: logger.WithFields(logrus.Fields{
			"struct": "DAO",
			"pkg": "dao",
		}),
	}
}

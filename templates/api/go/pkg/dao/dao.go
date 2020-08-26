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

// Converts a statement with 
func parameterizeStatement(query string) (string, []interface{}, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return "", nil, err
	}

	bindings := map[string]*querypb.BindVariable{}
	// dbcore is a random choice for a binding variable prefix so
	// we can be more certain when we find-replace on it
	sqlparser.Normalize(stmt, bindings, "dbcore")

	// Map :dbcore[0-9]* to $[0-9]*, will screw up if this literal
	// appeared in text in the query.
	re := regexp.MustCompile(":dbcore([0-9]*)")

	var invalidValue error
	var args []interface{}
	stmtWithBindings := re.ReplaceAllStringFunc(sqlparser.String(stmt), func (match string) string {
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

	return stmtWithBindings, args, invalidValue
}

func ParseFilter(filter string) (*Filter, error) {
	if filter == "" {
		return &Filter{}, nil
	}

	// TODO: validate filter uses acceptable subset of WHERE

	// Add stub select to make filter into a statement
	stmt, args, err := parameterizeStatement("SELECT x FROM x WHERE " + filter)
	if err != nil {
		return nil, err
	}

	return &Filter{
		// Take only the filter part from the statement
		filter: stmt[len("SELECT x FROM x "):],
		args: args,
	}, nil
}

// map $-prefixed variables from request context that can be turned
// into parameterized queries
func applyVariablesFromContext(filter string, ctx map[string]interface{}) string {
	re := regexp.MustCompile("\\$[a-zA-Z_]+")
	return re.ReplaceAllStringFunc(filter, func (match string) string {
		mapping, ok := ctx[match[1:]] // skip $ prefix
		if mapping == nil || !ok {
			return "NULL"
		}

		// Format as Go literal, probably works as a SQL literal too
		return fmt.Sprintf("%#v", mapping)
	})
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

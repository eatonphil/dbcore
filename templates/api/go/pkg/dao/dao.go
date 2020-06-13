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

type Filter struct {
	args []interface{}
	filter string
}

func ParseFilter(filter string) (*Filter, error) {
	if filter == "" {
		return &Filter{}, nil
	}

	stmt, err := sqlparser.Parse("SELECT 1 WHERE " + filter)
	if err != nil {
		return nil, err
	}

	bindings := map[string]*querypb.BindVariable{}
	// Extract out the literals into bind variables
	sqlparser.Normalize(stmt, bindings, "dbcore")

	// Map :dbcore[0-9]* to $[0-9]*
	exp := sqlparser.String(stmt.(*sqlparser.Select).Where.Expr)
	re := regexp.MustCompile(":dbcore([0-9]*)")

	var invalidValue error
	var f Filter
	f.filter = "WHERE " + re.ReplaceAllStringFunc(exp, func (match string) string {
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

			f.args = append(f.args, i)
		} else if v.IsFloat() {
			fl, err := strconv.ParseFloat(s, 64)
			if err != nil {
				invalidValue = err
				return ""
			}

			f.args = append(f.args, fl)
		} else if v.IsText() || v.IsQuoted() {
			f.args = append(f.args, s)
		} else if v.IsNull() {
			f.args = append(f.args, nil)
		} else {
			invalidValue = fmt.Errorf(`Unsupported value: "%s"`, s)
		}

		{{ if database.dialect == "postgres" }}
		nth, _ := strconv.ParseInt(re.FindStringSubmatch(exp)[1], 10, 64)
		return fmt.Sprintf("$%d", nth - 1)
		{{ else if database.dialect == "mysql" }}
		return "?"
		{{ end }}
	})

	if invalidValue != nil {
		return nil, invalidValue
	}

	return &f, nil
}

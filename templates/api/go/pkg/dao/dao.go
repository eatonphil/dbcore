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
		return &Filter{}, nil
	}

	stmt, err := sqlparser.Parse("SELECT 1 WHERE " + filter)
	if err != nil {
		return nil, err
	}

	bindings := map[string]*querypb.BindVariable{}
	// Extract out the literals into bind variables
	sqlparser.Normalize(stmt, bindings, "dbcore")
	exp := sqlparser.String(stmt.(*sqlparser.Select).Where.Expr)

	// Map :dbcore[0-9]* to $[0-9]*
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
		return fmt.Sprintf("$%d", len(f.args))
		{{ else if database.dialect == "mysql" || database.dialect == "sqlite" }}
		return "?"
		{{ end }}
	})

	if invalidValue != nil {
		return nil, invalidValue
	}

	return &f, nil
}

func ParseFilterWithContext(filter string, ctx map[string]interface{}) (*Filter, error) {
	re := regexp.MustCompile("\\$[a-zA-Z_]+")
	filter = re.ReplaceAllStringFunc(filter, func (match string) string {
		mapping, ok := ctx[match[1:]] // skip $ prefix
		if mapping == nil || !ok {
			return "NULL"
		}

		switch v := mapping.(type) {
		case bool:
			if v {
				return "TRUE"
			}
			
			return "FALSE"
		case string:
			return `"`+v+`"`
		default:
			return fmt.Sprintf("%v", v)
		}
	})

	return ParseFilter(filter)
}

type DAO struct {
	db *sqlx.DB
	logger logrus.FieldLogger
}

func (d DAO) IsAllowed(table, filter string, ctx map[string]interface{}) bool {
	sel := "SELECT COUNT(1)"
	from := fmt.Sprintf(`FROM "%s"`, table)

	// Allow override of select and from parts if specified
	stmt, err := sqlparser.Parse(filter)
	if err == nil {
		sel = sqlparser.String(stmt.(*sqlparser.Select).SelectExprs)
		from = sqlparser.String(stmt.(*sqlparser.Select).From)
		filter = sqlparser.String(stmt.(*sqlparser.Select).Where)[len("WHERE "):]
	}

	fmt.Println(sel, 1, from, 2, filter)

	f, err := ParseFilterWithContext(filter, ctx)
	if err != nil {
		d.logger.Warnf("Failed parsing allow filter: %s", err)
		return false
	}

	query := fmt.Sprintf(`%s %s %s`, sel, from, f.filter)
	d.logger.Debug(query)
	row := d.db.QueryRowx(query, f.args...)

	var count uint
	err = row.Scan(&count)
	if err != nil {
		d.logger.Warnf("Error fetching allow count: %s", err)
		return false
	}

	return count > 0
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

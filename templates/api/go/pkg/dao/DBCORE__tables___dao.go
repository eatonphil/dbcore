package dao

import (
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

{{~
func toGoType
  case $0.type
    when "int", "integer"
      if $0.nullable
        "null.Int"
      else
        "int32"
      end
    when "bigint"
      if $0.nullable
        "null.Int"
      else
        "int64"
      end
    when "text", "varchar", "char"
      if $0.nullable
        "null.String"
      else
        "string"
      end
    when "boolean"
      if $0.nullable
        "null.Bool"
      else
        "bool"
      end
    when "timestamp", "timestamp with time zone"
      if $0.nullable
        "null.Time"
      else
        "time.Time"
      end
    else
      "Unsupported type: " + $0.type
  end
end
~}}

type {{ table.label|dbcore_capitalize }} struct {
	{{~ for column in table.columns ~}}
	C_{{ column.name }} {{ toGoType column }} `db:"{{ column.name }}" json:"{{ column.name }}"`
	{{~ end ~}}
}

type {{ table.label|dbcore_capitalize }}PaginatedResponse struct {
	Total uint64 `json:"total"`
	Data []{{ table.label|dbcore_capitalize }} `json:"data"`
}

func (d DAO) {{ table.label|dbcore_capitalize }}GetMany(
	where *Filter,
	p Pagination,
	baseWhere string,
	baseCtx map[string]interface{},
) (*{{ table.label|dbcore_capitalize }}PaginatedResponse, error) {
	if where == nil {
		where = &Filter{}
{{ if api.audit.enabled && api.audit.deleted_at }}
		where.filter = `WHERE "{{ api.audit.deleted_at }}" IS NULL`
{{~ end ~}}
	} {{~ if api.audit.enabled && api.audit.deleted_at ~}} else {
	where.filter = where.filter + ` AND
  "{{ api.audit.deleted_at }}" IS NULL`
	}{{~ end ~}}

	{{~ if api.auth.enabled ~}}
	if baseWhere != "" {
		stmt, args, err := d.{{ table.label }}FilterToCompleteSQLStatement(baseWhere, baseCtx)
		if err != nil {
			// if baseWhere != "", this should only happen
			// during development with a bad filter
			panic(fmt.Errorf("Failed parsing base filter for get many request: %s", err))
		}

		// Combine base filter and where filter strings and args
		// TODO: handle restrictions on tables without a primary key
		where.filter = where.filter + ` AND "{{ table.primary_key.value.column }}" IN (` + stmt + ")"
		where.args = append(where.args, args...)
	}
	{{~ end ~}}

	query := fmt.Sprintf(`
SELECT
{{~ for column in table.columns ~}}
  "{{ column.name }}"{{ if !for.last || database.dialect != "sqlite" }},{{ end }}
{{~ end ~}}
{{~ if database.dialect != "sqlite" ~}}
  COUNT(1) OVER() AS __total
{{~ end ~}}
FROM
  "{{ table.name }}" t
%s
ORDER BY
  %s
LIMIT %d
OFFSET %d`, where.filter, p.Order, p.Limit, p.Offset)
	d.logger.Debug(query)
	rows, err := d.db.Queryx(query, where.args...)
	if err != nil {
		return nil, fmt.Errorf("Error in query: %s", err)
	}
	defer rows.Close()

	var response {{ table.label|dbcore_capitalize }}PaginatedResponse
	response.Data = []{{ table.label|dbcore_capitalize }}{}
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		var row struct {
			{{ table.label|dbcore_capitalize }}
{{~ if database.dialect != "sqlite" ~}}
			Total uint64 `db:"__total"`
{{~ end ~}}
		}
		err := rows.StructScan(&row)
		if err != nil {
			return nil, fmt.Errorf("Error scanning struct: %s", err)
		}

{{~ if database.dialect != "sqlite" ~}}
		response.Total = row.Total
{{~ end ~}}
		response.Data = append(response.Data, row.{{ table.label|dbcore_capitalize }})
	}

	{{~ if database.dialect == "sqlite" ~}}
	// COUNT() OVER() doesn't seem to work in the Go SQLite
	// package even though it works in the sqlite3 CLI.
	query = fmt.Sprintf(`
SELECT
  COUNT(1)
FROM
  "{{table.name}}"
%s
ORDER BY
  %s`, where.filter, p.Order)
	d.logger.Debug(query)
	row := d.db.QueryRowx(query, where.args...)
	err = row.Scan(&response.Total)
	if err != nil {
		return nil, fmt.Errorf("Error fetching total: %s", err, query)
	}
	{{~ end ~}}

	err = rows.Err()
	return &response, err
}

func (d DAO) {{ table.label|dbcore_capitalize }}Insert(body *{{ table.label|dbcore_capitalize }}) error {
	query := `
INSERT INTO {{ table.name }} (
{{~ for column in table.columns_no_audit ~}}
  "{{ column.name }}"{{ if !for.last }},{{ end }}
{{~ end ~}}
{{~ if api.audit.enabled ~}}
, "{{ api.audit.created_at }}", "{{ api.audit.updated_at }}"
{{~ end ~}})
VALUES (
{{~ for column in table.columns_no_audit ~}}
  {{ if database.dialect == "postgres" }}${{ for.index + 1 }}{{ else }}?{{ end }}{{ if !for.last }}, {{ end }}
{{~ end ~}},
{{~ if api.audit.enabled ~}}
{{~ if database.dialect == "sqlite" ~}}
  DATETIME('now'),
  DATETIME('now')
{{~ else ~}}
  NOW(),
  NOW()
{{~ end ~}}
{{~ end ~}})`
	d.logger.Debug(query)

	{{~ if database.dialect == "postgres" ~}}
	row := d.db.QueryRowx(query +`
RETURNING {{ if table.primary_key.value }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }}
`,
{{~ for column in table.columns_no_audit ~}}
		body.C_{{ column.name }},
{{~ end ~}}
	)
	return row.Scan(&body.C_{{ if table.primary_key.value }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }})
{{~ else if database.dialect == "mysql" || database.dialect == "sqlite" ~}}
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	{{ if database.dialect == "mysql" || database.dialect == "sqlite" }}var res sql.Result{{ end }}
	{{ if database.dialect == "mysql" || database.dialect == "sqlite" }}res{{ else }}_{{ end }}, err = stmt.Exec(
{{~ for column in table.columns_no_audit ~}}
		body.C_{{ column.name }},
{{~ end ~}})
	if err != nil {
		return err
	}

{{~ if table.primary_key.value ~}}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	body.C_{{ table.primary_key.value.column }} = {{ toGoType table.primary_key.value }}(id)
{{~ end ~}}
	return nil
{{~ end ~}}
}

{{ if table.primary_key.value }}
func (o {{ table.label|dbcore_capitalize }}) Id() {{ toGoType table.primary_key.value }} {
	return body.C_{{ table.primary_key.value.column }}
}

func (d DAO) {{ table.label|dbcore_capitalize }}Get(
	key {{ toGoType table.primary_key.value }},
) (*{{ table.label|dbcore_capitalize }}, error) {
	where, _ := ParseFilter(fmt.Sprintf("{{ table.primary_key.value.column }} = %#v", key))
	pagination := Pagination{
		Limit: 1,
		Offset: 0,
		Order: fmt.Sprintf("{{ table.primary_key.value.column }} DESC"),
	}

	r, err := d.{{ table.label|dbcore_capitalize }}GetMany(where, pagination, "", nil)
	if err != nil {
		return nil, err
	}

	if r.Total != 1 {
		return nil, ErrNotFound
	}

	return &r.Data[0], nil
}

func (d DAO) {{ table.label|dbcore_capitalize }}Update(key {{ toGoType table.primary_key.value }}, body {{ table.label|dbcore_capitalize }}) error {
	query := `
UPDATE
  "{{ table.name }}"
SET
{{~ for column in table.columns_no_audit ~}}
  "{{column.name}}" = {{ if database.dialect == "postgres" }}${{ index }}{{ else }}?{{ end }},
{{~ end ~}}
{{~ if database.dialect == "sqlite" ~}}
  "{{ api.audit.updated_at }}" = DATETIME('now')
{{~ else ~}}
  "{{ api.audit.updated_at }}" = NOW()
{{~ end ~}}
WHERE
{{~ if database.dialect == "postgres" ~}}
  "{{ table.primary_key.value.column }}" = ${{ table.columns_no_audit | array.size + 1 }}
{{~ else ~}}
  "{{ table.primary_key.value.column }}" = ?
{{~ end ~}}`
	d.logger.Debug(query)
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil
	}

	_, err = stmt.Exec(
{{~ for column in table.columns_no_audit ~}}
		body.C_{{ column.name }},
{{~ end ~}}
		key)
	return err
}

func (d DAO) {{ table.label|dbcore_capitalize }}Delete(key {{ toGoType table.primary_key.value }}) error {
	query := `
{{~ if api.audit.enabled && api.audit.deleted_at ~}}
UPDATE
  "{{ table.name }}"
SET
{{~ if database.dialect == "sqlite" }}
  "{{ api.audit.deleted_at }}" = DATETIME('now')
{{~ else ~}}
  "{{ api.audit.deleted_at }}" = NOW()
{{~ end ~}}
{{~ else ~}}
DELETE
  FROM "{{ table.name }}"
{{~ end ~}}
WHERE
  "{{ table.primary_key.value.column }}" = {{ if database.dialect == "postgres" }}$1{{ else }}?{{ end }}`

	d.logger.Debug(query)
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key)
	return err
}

func (d DAO) {{ table.label }}FilterToCompleteSQLStatement(
	filter string,
	ctx map[string]interface{},
) (string, []interface{}, error) {
	query := applyVariablesFromContext(filter, ctx)

	selectFromPrefix := ""
	// Allow override of select and from parts if specified
	_, err := sqlparser.Parse(query)
	if err != nil {
		// TODO: handle restrictions on tables without a primary key

		// TODO: Replace hack around using a mysql parser
		// where the table name must be backtick-quoted but
		// ANSI standard is double quotes
		selectFromPrefix = "SELECT \"{{ table.primary_key.value.column }}\" FROM `{{ table.name }}` WHERE "
		query = selectFromPrefix + query
	}

	parameterized, args, err := parameterizeStatement(query)
	if err != nil {
		return "", nil, fmt.Errorf("Failed parameterizing statement, `%s`: %s", query, err)
	}

	// -4 to get filter since the parser drops the quotes
	filter = parameterized[len(selectFromPrefix)-5:]
	// Replace backticks with quotes per ANSI standard
	return strings.ReplaceAll(selectFromPrefix, "`", `"`) + filter, args, nil
}

func (d DAO) {{ table.label|dbcore_capitalize }}IsAllowed(filter string, ctx map[string]interface{}) bool {
	query, args, err := d.{{ table.label }}FilterToCompleteSQLStatement(filter, ctx)
	if err != nil {
		d.logger.Warnf("Failed parsing allow filter: %s", err)
		return false
	}

	query = fmt.Sprintf(`SELECT COUNT(1) FROM (%s)`, query)
	d.logger.Debug(query)
	row := d.db.QueryRowx(query, args...)

	var count uint
	err = row.Scan(&count)
	if err != nil {
		d.logger.Warnf("Error fetching allow count: %s", err)
		return false
	}

	return count > 0
}

{{ end }}

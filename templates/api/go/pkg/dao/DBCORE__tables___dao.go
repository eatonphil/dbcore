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
        "Unsupported PostgreSQL type: " + $0.type
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

func (d DAO) {{ table.label|dbcore_capitalize }}GetMany(where *Filter, p Pagination) (*{{ table.label|dbcore_capitalize }}PaginatedResponse, error) {
	if where == nil {
		where = &Filter{}
	}{{~ if api.audit.deleted_at ~}} else {
	// Appended to an implicit `deleted_at IS NULL` base filter.
	where.filter = ` AND
  ` + where.filter[len("where "):]
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
{{~ if api.audit.deleted_at ~}}
WHERE
  "{{ api.audit.deleted_at }}" IS NULL{{~ end ~}} %s
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
	{{~ if api.audit.created_at ~}}
	{{~ if database.dialect == "sqlite" ~}}
	now := time.Now().Format(time.RFC3339)
	body.C_{{ api.audit.created_at }} = now
	{{~ else ~}}
	body.C_{{ api.audit.created_at }} = time.Now()
	{{~ end ~}}
	{{~ end ~}}

	{{~ if api.audit.updated_at ~}}
	{{~ if database.dialect == "sqlite" ~}}
	body.C_{{ api.audit.updated_at }} = now
	{{~ else ~}}
	body.C_{{ api.audit.updated_at }} = time.Now()
	{{~ end ~}}
	{{~ end ~}}

	{{~ if api.audit.deleted_at ~}}
	{{~ if database.dialect == "sqlite" ~}}
	body.C_{{ api.audit.updated_at }} = null.StringFromPtr(nil)
	{{~ else ~}}
	body.C_{{ api.audit.deleted_at }} = null.TimeFromPtr(nil)
	{{~ end ~}}
	{{~ end ~}}

	query := `
	INSERT INTO {{ table.name }} (
  {{~ for column in table.columns ~}}
  {{~ if column.auto_increment
         continue
        end ~}}
  "{{ column.name }}"{{ if !for.last }},{{ end }}
  {{~ end ~}})
VALUES (
  {{~ index = 0 ~}}
  {{~ for column in table.columns ~}}
  {{~ if column.auto_increment
         continue
      end ~}}
  {{ if database.dialect == "postgres" }}${{ index + 1 }}{{ else }}?{{ end }}{{ if !for.last }}, {{ end }}
  {{~ index = index + 1 ~}}
  {{~ end ~}})`
	d.logger.Debug(query)

	{{~ if database.dialect == "postgres" ~}}
	row := d.db.QueryRowx(query +`
RETURNING {{ if table.primary_key.value }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }}
`, {{~ for column in table.columns ~}}{{~ if column.auto_increment
		continue
		end ~}}body.C_{{ column.name }}{{ if !for.last }}, {{ end }}{{ end }})
	return row.Scan(&body.C_{{ if table.primary_key.value }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }})
	{{~ else if database.dialect == "mysql" || database.dialect == "sqlite" ~}}
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	{{ if database.dialect == "mysql" || database.dialect == "sqlite" }}var res sql.Result{{ end }}
	{{ if database.dialect == "mysql" || database.dialect == "sqlite" }}res{{ else }}_{{ end }}, err = stmt.Exec(
		{{~ for column in table.columns ~}}
		{{~ if column.auto_increment
		      continue
	            end ~}}
		body.C_{{ column.name }}{{ if !for.last }},{{ else }}){{ end }}{{ end }}
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
func (d DAO) {{ table.label|dbcore_capitalize }}Get(key {{ toGoType table.primary_key.value }}) (*{{ table.label|dbcore_capitalize }}, error) {
	where, _ := ParseFilter(fmt.Sprintf("{{ table.primary_key.value.column }} = %#v", key))
	pagination := Pagination{
		Limit: 1,
		Offset: 0,
		Order: fmt.Sprintf("{{ table.primary_key.value.column }} DESC"),
	}
	r, err := d.{{ table.label|dbcore_capitalize }}GetMany(where, pagination)
	if err != nil {
		return nil, err
	}

	if r.Total != 1 {
		return nil, ErrNotFound
	}

	return &r.Data[0], nil
}

func (d DAO) {{ table.label|dbcore_capitalize }}Update(key {{ toGoType table.primary_key.value }}, body {{ table.label|dbcore_capitalize }}) error {
	{{~ if api.audit.updated_at ~}}
	{{~ if database.dialect == "sqlite" ~}}
	body.C_{{ api.audit.updated_at }} = null.StringFromPtr(nil)
	{{~ else ~}}
	body.C_{{ api.audit.updated_at }} = time.Now()
	{{~ end ~}}
	{{~ end ~}}

	query := `
UPDATE
  "{{ table.name }}"
SET
  {{~ for column in table.columns ~}}
  "{{column.name}}" = {{ if database.dialect == "postgres" }}${{ for.index + 1 }}{{ else }}?{{ end }}{{ if !for.last }},{{ end }}
  {{~ end ~}}
WHERE
  {{ table.primary_key.value.column }} = {{ if database.dialect == "postgres" }}${{ table.columns | array.size + 1 }}{{ else }}?{{ end }}
`
	d.logger.Debug(query)
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil
	}

	_, err = stmt.Exec({{ for column in table.columns }}body.C_{{ column.name }}{{ if !for.last }},{{ end }}{{ end }}, key)
	return err
}

func (d DAO) {{ table.label|dbcore_capitalize }}Delete(key {{ toGoType table.primary_key.value }}) error {
	query := `
{{~ if api.audit.deleted_at ~}}
UPDATE
  "{{ table.name }}"
SET "{{ api.audit.deleted_at }}" = {{ database.dialect == "sqlite" }}DATETIME('now'){{ else }}NOW(){{ end }}
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
{{ end }}

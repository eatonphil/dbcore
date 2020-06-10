package dao

import (
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

type {{ table.name|string.capitalize }} struct {
	{{~ for column in table.columns ~}}
	C_{{ column.name }} {{ column.go_type }} `db:"{{ column.name }}" json:"{{ column.name }}"`
	{{~ end ~}}
}

type {{ table.name|string.capitalize }}PaginatedResponse struct {
	Total uint64 `json:"total"`
	Data []{{ table.name|string.capitalize }} `json:"data"`
}

func (d DAO) {{ table.name|string.capitalize }}GetMany(where *Filter, p Pagination) (*{{ table.name|string.capitalize }}PaginatedResponse, error) {
	if where == nil {
		where = &Filter{}
	}

	query := fmt.Sprintf(`
SELECT
  {{~ for column in table.columns ~}}
  "{{ column.name }}",
  {{~ end ~}}
  COUNT(1) OVER () AS __total
FROM
  "{{table.name}}"
%s
ORDER BY
  %s
LIMIT %d
OFFSET %d`, where.filter, p.Order, p.Limit, p.Offset)
	d.logger.Debug(query)
	rows, err := d.db.Queryx(query, where.args...)
	if err != nil {
		return nil, err
	}

	var response {{ table.name|string.capitalize }}PaginatedResponse
	response.Data = []{{ table.name|string.capitalize }}{}
	for rows.Next() {
		var row struct {
			{{ table.name|string.capitalize }}
			Total uint64 `db:"__total"`
		}
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		response.Total = row.Total
		response.Data = append(response.Data, row.{{ table.name|string.capitalize }})
	}

	return &response, err
}

func (d DAO) {{ table.name|string.capitalize }}Insert(body *{{ table.name|string.capitalize }}) error {
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
RETURNING {{ if table.primary_key.is_some }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }}
`, {{~ for column in table.columns ~}}{{~ if column.auto_increment
		continue
		end ~}}body.C_{{ column.name }}{{ if !for.last }}, {{ end }}{{ end }})
	return row.Scan(&body.C_{{ if table.primary_key.is_some }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }})
	{{~ else if database.dialect == "mysql" ~}}
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	{{ if database.dialect == "postgres" }}var res sql.Result{{ end }}
	{{ if database.dialect == "postgres" }}res{{ else }}_{{ end }}, err = stmt.Exec(
		{{~ for column in table.columns ~}}
		{{~ if column.auto_increment
		      continue
	            end ~}}
		body.C_{{ column.name }}{{ if !for.last }},{{ else }}){{ end }}{{ end }}
	if err != nil {
		return err
	}

	{{~ if table.primary_key.is_some ~}}
	body.C_{{ column.primary_key.value.column }}, err = res.LastInsertId()
	if err != nil {
		return err
	}
	{{~ end ~}}
	return nil
	{{~ end ~}}
}

{{ if table.primary_key.is_some }}
func (d DAO) {{ table.name|string.capitalize }}Update(key {{ table.primary_key.go_type }}, body {{ table.name|string.capitalize }}) err {
	stmt, err := d.db.Prepare(`
UPDATE
  "{{ table.name }}"
SET
  {{~ for column in table.columns ~}}
  "{{column.name}}" = {{ if database.dialect == "postgres" }}${{ for.index + 2 }}{{ else }}?{{ end }}{{ if !for.last }},{{ end }}
  {{~ end ~}}
WHERE
  {{ table.primary_key.name }} = $1
`)
	if err != nil {
		return nil
	}

	_, err = stmt.Exec(id, {{ for column in table.columns }}body.C_{{ column.name }}{{ if !for.last }},{{ end }}{{ end }})
	return err
}

func (d DAO) {{ table.name|string.capitalize }}Delete(key {{ table.primary_key.go_type }}) error {
	stmt, err := d.db.Prepare(`
DELETE
  FROM "{{ table.name }}"
WHERE
  "{{ table.primary_key.value.column }}" = {{ if database.dialect == "postgres" }}$1{{ else }}?{{ end }}`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key)
	return err
}
{{ end }}

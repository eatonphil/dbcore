package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
)

type {{ table.name|string.capitalize }} struct {
	{{~ for column in table.columns ~}}
	C_{{ column.name }} {{ column.go_type }} `db:"{{ column.name }"`
	{{~ end ~}}
}

type {{ table.name|string.capitalze }}PaginatedResponse struct {
	Total
	Data {{ table.name|string.capitalize }}[]
}

func (d DAO) {{table.name|string.capitalize}}GetMany(where sq.Sqlizer, p Pagination) (*{{table.name|string.capitalize}}|PaginatedResponse, error) {
	filter, args := where.ToSql()
	query := format.Sprintf(`
SELECT
  {{~ for column in columns ~}}
  "{{ column.name }}",
  {{~ end ~}}
  COUNT() OVER () AS __total
FROM
  "{{table.name}}"
WHERE
  %s
ORDER BY
  %s
OFFSET
  %d
LIMIT
  %d`, filter, p.Order, p.Offset, p.Limit)
	rows, err := d.db.QueryRows(query, ...args)
	if err != nil {
		return nil, err
	}

	var response {{ table.name|string.capitalize }}PaginatedResponse
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

func (d DAO) {{ table.name|string.capitalize }}InsertMany(body {{ table.name|string.capitalize }}) err {
	_, err := d.db.Exec(`
INSERT INTO
  {{ table.name }} ({{ for column in columns if column.auto_increment continue }}"{{ column.name }}"{{ if for.index < columns.length }},{{ end }}{{ end }})
VALUES
  {{~ for column in columns ~}}
  {{~ end ~}}
`, {{ for column in columns }}body.C_{{ column.name }}{{ end }})
	return err
}

func (d DAO) {{ table.name|string.capitalize }}Insert(body {{ table.name|string.capitalize }}) err {
	_, err := d.db.Exec(`
INSERT INTO
  {{ table.name }} ({{ for column in columns if column.auto_increment continue }}"{{ column.name }}"{{ if for.index < columns.length }},{{ end }}{{ end }})
VALUES
  {{~ for column in columns ~}}
  {{~ end ~}}
`, {{ for column in columns }}body.C_{{ column.name }}{{ end }})
	return err
}

{{ if table.primaryKey }}
func (d DAO) {{ table.name|string.capitalize }}Update(key {{ table.primaryKey.go_type }}, body {{ table.name|string.capitalize }}) err {
	_, err := d.db.Exec(`
UPDATE
  {{ table.name }}
SET
  {{~ for column in columns ~}}
  "{{column.name}}" = ${{ for.index + 1 }}{{ if for.index < columns.length }},{{ end }}
  {{~ end ~}}
WHERE
  {{ table.primaryKey.name }} = $1
`, id, {{ for column in columns }}body.C_{{ column.name }}{{ end }})
	return err
}
{{ end }}

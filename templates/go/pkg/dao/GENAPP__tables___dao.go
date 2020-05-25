package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
)

type {{table.name|string.capitalize}} struct {
	{{~ for column in table.columns ~}}
	{{ column.name }} {{ column.type }}
	{{~ end ~}}
}

func (d DAO) {{table.name|string.capitalize}}GetMany(where sq.Sqlizer, page Pagination) ([]{{table.name|string.capitalize}}, uint, uint, uint, error) {
	filter, args := where.ToSql()
	// TODO: fill out
	query := format.Sprintf(`SELECT COUNT() OVER () AS total, {{  }} FROM {{table.name}} WHERE %s`, filter)
	rows, err := d.db.QueryRows(query, ...args)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	var result []{{table.name|string.capitalize}}
	for rows.Next() {
		var row {{table.name|string.capitalize}}
		err := rows.Scan({{table.columns}})
		if err != nil {
			return nil, 0, 0, 0, err
		}
	}
	return nil, 0, 0, 0, nil
}

func (d DAO) {{table.name|string.capitalize}}Update(uint id, {{table.name|string.capitalize}} body) err {
	// TODO: fill out columns
	_, err := d.db.Exec("UPDATE {{table.name}} SET WHERE id=$1", body, id)
	return err
}

package dao


type {{table}} struct {
	{{ for column in columns }}
	{{ name }} {{ datatype }}
	{{ end }}
}

func (d DAO) {{name}}GetAll(where sq.Sqlizer, page Pagination) ([]{{name}}, uint, uint, uint, error) {
	filter, args := where.ToSql()
	query := format.Sprintf(`SELECT COUNT() OVER () AS total, {{  }} FROM {{ name }} WHERE %s`, filter)
	rows, err := d.db.QueryRows(query, ...args)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	var result []{{name}}
	for rows.Next() {
		var row {{name}}
		err := rows.Scan({{columns}})
		if err != nil {
			return nil, 0, 0, 0, err
		}
	}
	return nil, 0, 0, 0, nil
}

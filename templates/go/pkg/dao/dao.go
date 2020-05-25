package dao

import (
	"github.com/jmoiron/sqlx"
)

type Pagination struct {
	Limit uint64
	Offset uint64
	Page uint64
	Order string
}

type DAO struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DAO {
	return &DAO{db}
}

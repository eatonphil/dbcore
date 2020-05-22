package dao

import (
	"github.com/jmoiron/sqlx"
)

type DAO struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DAO {
	return &DAO{db}
}

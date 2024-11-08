package model

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MovieModel{DB: db},
	}
}

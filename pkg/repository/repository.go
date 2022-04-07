package repository

import "github.com/jmoiron/sqlx"

type Authentication interface {

}

type User interface {

}

type Repository struct {

}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}

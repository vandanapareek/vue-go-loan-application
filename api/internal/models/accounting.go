package models

import (
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

type ProviderRepo interface {
	GetAllProviders() ([]*Provider, error)
}

type providerRepo struct {
	db *sqlx.DB
}

func NewProviderRepo(db *sqlx.DB) ProviderRepo {
	return &providerRepo{
		db: db,
	}
}

func (r *providerRepo) GetAllProviders() ([]*Provider, error) {
	query := `SELECT
			name,
			slug
		FROM providers
		WHERE status = 1`
	query = r.db.Rebind(query)
	var userAcc []*Provider
	err := r.db.Select(&userAcc, query)
	return userAcc, err
}

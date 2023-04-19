package models

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Username string `json:"username"`
	Currency string `json:"currency"`
}

type UserRepo interface {
	GetByUserName(username string) (*User, error)
	GetUserInfo(username string) (*UserInfo, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetByUserName(username string) (*User, error) {
	query := `SELECT
			username,
			password
		FROM users
		WHERE username = ?`

	query = r.db.Rebind(query)
	var user User
	err := r.db.Get(&user, query, username)
	return &user, err
}

func (r *userRepo) GetUserInfo(username string) (*UserInfo, error) {
	query := `SELECT
			username
		FROM users
		WHERE username = ?  LIMIT 1`

	query = r.db.Rebind(query)
	var userAcc UserInfo
	err := r.db.Get(&userAcc, query, username)
	return &userAcc, err
}

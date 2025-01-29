package data

import (
	"database/sql"
	"log/slog"

	"github.com/yunya101/perepisochnik/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Insert(username string) {
	stmt := `INSERT INTO users (username)
			VALUES ($1)`

	_, err := r.DB.Exec(stmt, username)

	if err != nil {
		slog.Error(err.Error())
	}

}

func (r *UserRepo) GetByName(username string) *models.User {
	stmt := `SELECT * FROM users
			WHERE username = $1`

	row := r.DB.QueryRow(stmt, username)

	user := &models.User{}

	row.Scan(&user.Username)

	return user
}

package data

import (
	"database/sql"

	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Insert(username, pass string) error {
	stmt := `INSERT INTO users (username, password)
			VALUES ($1, $2)`

	_, err := r.DB.Exec(stmt, username, pass)

	if err != nil {
		conf.ErrLog.Printf("%s:%s", err.Error(), username)
		return err
	}

	conf.InfoLog.Printf("Insert into users:%s", username)
	return nil

}

func (r *UserRepo) GetByName(username string) (*models.User, error) {
	stmt := `SELECT * FROM users
			WHERE username = $1`

	row := r.DB.QueryRow(stmt, username)

	user := &models.User{}

	if err := row.Scan(&user.Username, &user.Password); err != nil {
		conf.ErrLog.Printf("%s:%s", err.Error(), username)
		return nil, err
	}

	conf.InfoLog.Printf("Select form users:%s", username)

	return user, nil
}

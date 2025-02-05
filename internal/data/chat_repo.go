package data

import (
	"database/sql"

	"github.com/yunya101/perepisochnik/internal/models"
)

type ChatRepo struct {
	db *sql.DB
}

func (r *ChatRepo) Insert(chat *models.Chat) {
	// stmt := `INSERT INTO user_chat ()`
}

func (r *ChatRepo) Update() {

}

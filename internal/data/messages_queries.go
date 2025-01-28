package data

import (
	"database/sql"
	"log/slog"

	"github.com/yunya101/perepisochnik/internal/models"
)

type MessageRepo struct {
	DB *sql.DB
}

func (r *MessageRepo) Insert(msg *models.Message) {
	stmt := `INSERT INTO messages (reciver, recipient, text, chat_id)
			VALUES ($1, $2, $3, $4)`

	_, err := r.DB.Exec(stmt, msg.Reciver, msg.Recipient, msg.Text, msg.Chat)

	if err != nil {
		slog.Error(err.Error())
	}
}

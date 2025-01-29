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
	stmt := `INSERT INTO messages (reciver, recipient, text)
			VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(stmt, msg.Reciver, msg.Recipient, msg.Text)

	if err != nil {
		slog.Error(err.Error())
	}
}

func (r *MessageRepo) GetAllByUsername(username string) []*models.Message {
	stmt := `SELECT * FROM messages
			WHERE reciver = $1 or recipient = $1`
	rows, err := r.DB.Query(stmt, username)

	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	defer rows.Close()

	result := []*models.Message{}

	for rows.Next() {
		msg := &models.Message{}
		rows.Scan(&msg.Id, &msg.Reciver, &msg.Recipient, &msg.Text)
		result = append(result, msg)
	}

	if rows.Err() != nil {
		slog.Error(rows.Err().Error())
		return nil
	}

	return result

}

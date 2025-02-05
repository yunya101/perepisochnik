package data

import (
	"database/sql"

	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/models"
)

type MessageRepo struct {
	DB *sql.DB
}

func (r *MessageRepo) Insert(msg *models.Message) error {
	stmt := `INSERT INTO messages (reciver, recipient, text, chat_id)
			VALUES ($1, $2, $3, $4)`

	_, err := r.DB.Exec(stmt, msg.Reciver, msg.Recipient, msg.Text, msg.ChatId)
	if err != nil {
		conf.ErrLog.Printf("%s:%v", err.Error(), msg)
		return err
	}

	conf.InfoLog.Printf("Insert into messages:%v", msg)
	return nil
}

func (r *MessageRepo) GetAllByUsername(username string) ([]*models.Message, error) {
	stmt := `SELECT * FROM messages
			WHERE reciver = $1 or recipient = $1`
	rows, err := r.DB.Query(stmt, username)

	if err != nil {
		conf.ErrLog.Printf("%s:%v", err.Error(), username)
		return nil, err
	}

	defer rows.Close()

	result := []*models.Message{}

	for rows.Next() {
		msg := models.Message{}
		if err := rows.Scan(&msg.Id, &msg.Reciver, &msg.Recipient, &msg.Text, &msg.ChatId); err != nil {
			conf.ErrLog.Printf("%s:%s", err.Error(), username)
			return nil, err
		} else {
			result = append(result, &msg)
		}
	}

	if rows.Err() != nil {
		conf.ErrLog.Println(err.Error())
		return nil, err
	}

	return result, nil

}

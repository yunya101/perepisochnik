package data

import (
	"database/sql"

	"github.com/lib/pq"

	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/models"
)

type ChatRepo struct {
	DB *sql.DB
}

func (r *ChatRepo) Insert(chat *models.Chat) error {

	stmt := `with new_chat AS (
			INSERT INTO chats (name)
			VALUES ($1)
			RETURNING id
		)

		INSERT INTO users_chats (username, chat)
		SELECT username, new_chat.id
		FROM unnest($2::text[]) AS username
		CROSS JOIN new_chat;`

	_, err := r.DB.Exec(stmt, chat.Name, pq.Array(chat.Users))

	if err != nil {
		conf.ErrLog.Printf("%s:%v", err.Error(), chat)
		return err
	}

	conf.InfoLog.Printf("Insert into chats successfully:%v", chat)

	return nil
}

func (r *ChatRepo) Update() {

}

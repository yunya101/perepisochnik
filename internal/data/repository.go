package data

import (
	"database/sql"

	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/models"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) SetDB(db *sql.DB) {
	r.db = db
}

func (r *Repository) GetChatsByUsername(username string) ([]*models.Chat, error) {
	stmt := `SELECT c.id, c.name FROM chats c
			JOIN users_chats uc ON c.id = uc.chat
			WHERE uc.username = $1`

	rows, err := r.db.Query(stmt, username)

	if err != nil {
		conf.ErrLog.Printf("%s:%s", err, username)
		return nil, err
	}

	defer rows.Close()

	var chats = make([]*models.Chat, 0)

	for rows.Next() {
		var chat *models.Chat = &models.Chat{}

		if err := rows.Scan(&chat.ID, &chat.Name); err != nil {
			conf.ErrLog.Printf("%s:%s", err.Error(), username)
			return nil, err
		}

		chats = append(chats, chat)

	}

	if rows.Err() != nil {
		conf.ErrLog.Printf("%s:%s", err, username)
		return nil, err
	}

	conf.InfoLog.Printf("User's chats was selected:%s", username)

	return chats, nil
}

func (r *Repository) GetUsersFromChat(chat *models.Chat) (*models.Chat, error) {

	stmt := `SELECT username FROM users_chats
			WHERE chat = $1`

	rows, err := r.db.Query(stmt, chat.ID)

	if err != nil {
		conf.ErrLog.Printf("%s:%v", err, chat.ID)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			conf.ErrLog.Printf("%s:%v", err, chat.ID)
			return nil, err
		}

		chat.Users = append(chat.Users, username)
	}

	if rows.Err() != nil {
		conf.ErrLog.Printf("%s:%v", err, chat.ID)
		return nil, err
	}

	conf.InfoLog.Printf("Chat's users selected:%v", chat.ID)

	return chat, nil

}

func (r *Repository) GetMsgsFromChat(chat *models.Chat) (*models.Chat, error) {
	stmt := `SELECT * FROM messages WHERE chat_id = $1`

	rows, err := r.db.Query(stmt, chat.ID)

	if err != nil {
		conf.ErrLog.Printf("%s:%v", err, chat.ID)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var msg models.Message

		if err := rows.Scan(&msg.Id, &msg.Sender, &msg.Text, &msg.ChatId); err != nil {
			conf.ErrLog.Printf("%s:%v", err, chat.ID)
			return nil, err
		}

		chat.Messages = append(chat.Messages, msg)
	}

	if rows.Err() != nil {
		conf.ErrLog.Printf("%s:%v", err, chat.ID)
		return nil, err
	}

	conf.InfoLog.Printf("Chat's messages selected:%v", chat.ID)
	return chat, nil
}

func (r *Repository) InsertMsg(msg *models.Message) error {
	stmt := `INSERT INTO messages (sender, text, chat_id)
	VALUES ($1, $2, $3)`

	_, err := r.db.Exec(stmt, msg.Sender, msg.Text, msg.ChatId)

	if err != nil {
		conf.ErrLog.Printf("%s:%v", err, msg)
		return err
	}

	return nil
}

package models

type Message struct {
	Id     int      `json:"id"`
	Sender string   `json:"sender"`
	Text   string   `json:"text"`
	ChatId int      `json:"chat_id"`
	ReadBy []string `json:"read_by"`
}

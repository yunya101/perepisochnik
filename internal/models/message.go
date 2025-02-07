package models

type Message struct {
	Id     int64  `json:"id"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	ChatId int64  `json:"chat_id"`
}

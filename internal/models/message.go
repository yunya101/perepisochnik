package models

type Message struct {
	Id        int    `json:"id"`
	Reciver   string `json:"reciver"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

package models

type Message struct {
	Id        int    `json:"id"`
	Reciver   int    `json:"reciver"`
	Recipient int    `json:"recipient"`
	Text      string `json:"text"`
}

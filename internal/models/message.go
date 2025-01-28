package models

type Message struct {
	Id        int    `json:"id"`
	Reciver   string `json:"reciver"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
	Chat      int    `json:"chat_id"`
}

type Chat struct {
	Id       int       `json:"id"`
	User1    User      `json:"user1"`
	User2    User      `json:"user2"`
	Messages []Message `json:"messages"`
}

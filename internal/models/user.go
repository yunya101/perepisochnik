package models

type User struct {
	Username string  `json:"username"`
	Chats    []*Chat `json:"chats"`
}

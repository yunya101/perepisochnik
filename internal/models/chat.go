package models

type Chat struct {
	ID       int       `json:"id"`
	Users    []User    `json:"users"`
	Messages []Message `json:"messages"`
}

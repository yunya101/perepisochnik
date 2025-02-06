package models

type Chat struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Users    []string  `json:"users"`
	Messages []Message `json:"messages"`
}

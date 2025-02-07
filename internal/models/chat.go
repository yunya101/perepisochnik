package models

type Chat struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Users    []string  `json:"users"`
	Messages []Message `json:"messages"`
}

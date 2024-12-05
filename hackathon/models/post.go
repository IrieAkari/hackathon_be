package models

type Post struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Content string `json:"content"`
}

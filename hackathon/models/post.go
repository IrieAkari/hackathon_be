package models

type Post struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Content string `json:"content"`
}

type PostWithUserName struct {
	Id        string  `json:"id"`
	UserId    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	ParentId  *string `json:"parent_id"`
}

package models

type UserResForHTTPGet struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserReqForHTTPPost struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

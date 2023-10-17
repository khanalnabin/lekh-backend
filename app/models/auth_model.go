package models

type Register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

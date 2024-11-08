package main

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PasswordHash []byte `json:"password"`
	Email        string `json:"email"`
}

type Account struct {
	UserID  int `json:"user_id"`
	Number  int `json:"number"`
	Balance int `json:"balance"`
}

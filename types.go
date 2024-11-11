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

type UserCreateRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PasswordHash []byte `json:"password"`
	Email        string `json:"email"`
}

func NewUser(u *UserCreateRequest) *User {
	return &User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
	}
}

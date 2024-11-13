package main

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
}

type Account struct {
	UserID  int `json:"user_id"`
	Number  int `json:"number"`
	Balance int `json:"balance"`
}

type RegisterUserRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

func NewUser(firstName, lastName, password, email string) (*User, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: string(pwd),
		Email:        email,
	}, nil
}

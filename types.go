package main

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Account
}

type Account struct {
	Number  int `json:"number"`
	Balance int `json:"balance"`
}

type RegisterUserRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UsersResponse struct {
	Users []*User `json:"users"`
}

func NewUser(firstName, lastName, email, password string) (*User, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(pwd),
	}, nil
}

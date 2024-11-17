package main

import (
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phoneNumber"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
	Number       string    `json:"number"`
	Balance      int       `json:"balance"`
}

type RegisterUserRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	PasswordHash string `json:"password_hash"`
}

type UsersResponse struct {
	Users []*User `json:"users"`
}

func NewUser(firstName, lastName, email, phoneNumber, password string) (*User, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PhoneNumber:  phoneNumber,
		PasswordHash: string(pwd),
		CreatedAt:    time.Now(),
		Number:       generateAccountNumber(),
		Balance:      0.00,
	}, nil
}

func generateAccountNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	number := ""

	for i := 1; i <= 19; i++ {
		if i == 5 || i == 10 || i == 15 {
			number += " "
			continue
		}

		number += strconv.Itoa(r.Intn(10))
	}

	return number
}

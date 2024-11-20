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
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	Number       string    `json:"number"`
	Balance      float64   `json:"balance"`
}

type RegisterUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type ChargeRequest struct {
	AccountNumber string  `json:"accountNumber"`
	Amount        float64 `json:"amount"`
}

type TransferRequest struct {
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	Amount      int    `json:"amount"`
}

type UsersResponse struct {
	Users []*User `json:"users"`
}

func NewUser(firstName, lastName, email, phoneNumber, password string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PhoneNumber:  phoneNumber,
		PasswordHash: string(passwordHash),
	}, nil
}

func generateAccountNumber() string {
	var (
		number = ""
		r      = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 1; i <= 19; i++ {
		if i%5 == 0 {
			number += " "
			continue
		}

		number += strconv.Itoa(r.Intn(10))
	}

	return number
}

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
	Number       string    `json:"number"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"createdAt"`
}

type RegisterUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type RegisterUserResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	User       `json:"user"`
}

type ChargeRequest struct {
	AccountNumber string  `json:"accountNumber"`
	Amount        float64 `json:"amount"`
}

type ChargeResponse struct {
	StatusCode int     `json:"statusCode"`
	Msg        string  `json:"msg"`
	Amount     float64 `json:"amount"`
	Balance    float64 `json:"balance"`
}

type TransferRequest struct {
	FromAccount string  `json:"fromAccount"`
	ToAccount   string  `json:"toAccount"`
	Amount      float64 `json:"amount"`
}

type TransferResponse struct {
	StatusCode int     `json:"statusCode"`
	Msg        string  `json:"msg"`
	Amount     float64 `json:"amount"`
	Balance    float64 `json:"balance"`
}

type DeleteUserResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	ID         int    `json:"id"`
}

type UserResponse struct {
	StatusCode int  `json:"statusCode"`
	User       User `json:"user"`
}

type UsersResponse struct {
	StatusCode int     `json:"statusCode"`
	Users      []*User `json:"users"`
}

func NewUser(firstName, lastName, email, phoneNumber, password string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	accountNumber := generateAccountNumber()

	return &User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PhoneNumber:  phoneNumber,
		CreatedAt:    time.Now(),
		PasswordHash: string(passwordHash),
		Number:       accountNumber,
		Balance:      0.00,
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

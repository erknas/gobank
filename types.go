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
	Acc          Account   `json:"account"`
}

type Account struct {
	ID      int     `json:"id"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
}

type UserResponse struct {
	StatusCode int  `json:"statusCode"`
	User       User `json:"user"`
}

type UsersResponse struct {
	StatusCode int    `json:"statusCode"`
	Users      []User `json:"users"`
}

type NewUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type NewUserResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	User       User   `json:"user"`
}

type DeleteUserResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	ID         int    `json:"id"`
}

type Transaction struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"-"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	FromAccount string    `json:"fromAccount,omitempty"`
	ToAccount   string    `json:"toAccount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TransactionRequest struct {
	Type        string  `json:"type"`
	FromAccount string  `json:"fromAccount"`
	ToAccount   string  `json:"toAccount"`
	Amount      float64 `json:"amount"`
}

type TransactionResponse struct {
	StatusCode  int         `json:"statusCode"`
	Msg         string      `json:"msg"`
	Transaction Transaction `json:"transaction"`
}

type TransactionsResponse struct {
	StatusCode   int           `json:"statusCode"`
	AccountID    int           `json:"accountId"`
	Transactions []Transaction `json:"transactions"`
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
		Acc: Account{
			Number:  accountNumber,
			Balance: 0,
		},
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

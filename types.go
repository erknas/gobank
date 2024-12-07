package main

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

const layout = "01/06"

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	PhoneNumber  string    `json:"phoneNumber"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	Account      Account   `json:"account"`
}

type Account struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
	Card    Card    `json:"card"`
}

type Card struct {
	ID         int    `json:"id"`
	Number     string `json:"number"`
	CVV        string `json:"-"`
	ExpireTime string `json:"expireTime"`
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
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type NewUserResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

type DeleteUserResponse struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	ID         int    `json:"id"`
}

type Transaction struct {
	ID             uuid.UUID `json:"id"`
	AccountID      int       `json:"-"`
	Type           string    `json:"type"`
	Amount         float64   `json:"amount"`
	FromCardNumber string    `json:"fromCardNumber,omitempty"`
	ToCardNumber   string    `json:"toCardNumber"`
	CreatedAt      time.Time `json:"createdAt"`
}

type TransactionRequest struct {
	Type           string  `json:"type"`
	FromCardNumber string  `json:"fromCardNumber"`
	ToCardNumber   string  `json:"toCardNumber"`
	Amount         float64 `json:"amount"`
}

type TransactionResponse struct {
	StatusCode  int         `json:"statusCode"`
	Msg         string      `json:"msg"`
	Transaction Transaction `json:"transaction"`
}

type TransactionsResponse struct {
	StatusCode   int           `json:"statusCode"`
	UserID       int           `json:"userId"`
	Transactions []Transaction `json:"transactions"`
}

func NewUser(firstName, lastName, phoneNumber, password string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	card := NewCard()

	return &User{
		FirstName:    firstName,
		LastName:     lastName,
		PhoneNumber:  phoneNumber,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now().UTC(),
		Account: Account{
			Balance: 0.00,
			Card: Card{
				Number:     card.Number,
				CVV:        card.CVV,
				ExpireTime: card.ExpireTime,
			},
		},
	}, nil
}

func NewCard() Card {
	var (
		cardNumber = ""
		cvv        = ""
		r          = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	)

	for i := 1; i <= 16; i++ {
		cardNumber += strconv.Itoa(r.Intn(10))
	}

	for i := 1; i <= 3; i++ {
		cvv += strconv.Itoa(r.Intn(10))
	}

	now := time.Now()
	expireTime := now.AddDate(0, 61, 0).Format(layout)

	return Card{
		Number:     cardNumber,
		CVV:        cvv,
		ExpireTime: expireTime,
	}
}

package main

import "context"

type Storer interface {
	Register(context.Context, *User) error
	Deposit(context.Context, *TransactionRequest) (Transaction, error)
	Transfer(context.Context, *TransactionRequest) (Transaction, error)
	UserByID(context.Context, int) (User, error)
	TransactionsByUser(context.Context, int) ([]Transaction, error)
	DeleteUser(context.Context, int) error
	Users(context.Context) ([]User, error)
}

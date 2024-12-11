package main

import "context"

type Storer interface {
	Register(context.Context, *User) (int, error)
	Deposit(context.Context, *TransactionRequest) (Transaction, error)
	Transfer(context.Context, *TransactionRequest) (Transaction, error)
	UserByID(context.Context, int) (User, error)
	TransactionsByUser(context.Context, int) ([]Transaction, error)
	Users(context.Context) ([]User, error)
}

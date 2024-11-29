package main

import "context"

type Storer interface {
	Register(context.Context, *User) (int, error)
	Charge(context.Context, *TransactionRequest) (*Transaction, error)
	Transfer(context.Context, *TransactionRequest) (*Transaction, error)
	GetUserByID(context.Context, int) (*User, error)
	Delete(context.Context, int) error
	GetUsers(context.Context) ([]*User, error)
	GetTransactions(context.Context) ([]*Transaction, error)
}

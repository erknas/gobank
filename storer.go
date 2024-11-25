package main

import "context"

type Storer interface {
	Register(context.Context, *User) (int, error)
	Charge(context.Context, *ChargeRequest) (float64, error)
	Transfer(context.Context, *TransferRequest) (float64, error)
	GetUserByID(context.Context, int) (*User, error)
	Delete(context.Context, int) error
	GetUsers(context.Context) ([]*User, error)
}

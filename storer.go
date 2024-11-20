package main

import "context"

type Storer interface {
	Register(context.Context, *User) error
	Charge(context.Context, *ChargeRequest) error
	Transfer(context.Context, *TransferRequest) error
	GetUserByID(context.Context, int) (*User, error)
	GetUsers(context.Context) ([]*User, error)
}

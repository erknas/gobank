package main

import "context"

type Storer interface {
	Register(context.Context, *User) error
	GetUserByID(context.Context, int) (*User, error)
	GetUsers(context.Context) ([]*User, error)
	Transfer(context.Context, string, string, int) error
}

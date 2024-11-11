package main

import "context"

type Storer interface {
	Register(context.Context, *UserCreateRequest) error
	Get(context.Context, string) (*User, error)
	Transfer(context.Context, string, string, int) error
}

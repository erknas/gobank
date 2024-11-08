package main

import "context"

type Storer interface {
	Register(context.Context, *Account) error
	Get(context.Context, string, []byte) (*Account, error)
	Transfer(context.Context, string, string, int) error
}

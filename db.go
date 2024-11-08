package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %s", err)
	}

	return &Storage{
		conn: conn,
	}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Register(ctx context.Context, acc *Account) error {
	return nil
}

func (s *Storage) Get(ctx context.Context, email string, pw []byte) (*Account, error) {
	return nil, nil
}

func (s *Storage) Transfer(ctx context.Context, from, to string, amount int) error {
	return nil
}

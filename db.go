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

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Storage{
		conn: conn,
	}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Register(ctx context.Context, user *UserCreateRequest) error {
	query := `INSERT INTO users(first_name, last_name, password_hash, email) VALUES ($1, $2, $3, $4)`
	_, err := s.conn.Exec(ctx, query, user.FirstName, user.LastName, user.PasswordHash, user.Email)

	return err
}

func (s *Storage) Get(ctx context.Context, name string) (*User, error) {
	query := `SELECT * FROM users where first_name=$1`

	rows, err := s.conn.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %s", name)
	}
	defer rows.Close()

	user := new(User)

	for rows.Next() {
		if err := rows.Scan(user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *Storage) Transfer(ctx context.Context, from, to string, amount int) error {
	return nil
}

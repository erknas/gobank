package main

import (
	"context"
	"fmt"
	"log"

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

func (s *Storage) Register(ctx context.Context, user *User) (int, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	defer func() { wrapErr(ctx, tx, err) }()

	var userID int
	if err = tx.QueryRow(ctx, createUserQuery, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.PasswordHash).Scan(&userID); err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, createAccountQuery, userID, user.Number)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*User, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})

	defer func() { wrapErr(ctx, tx, err) }()

	user := new(User)

	err = s.conn.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Number, &user.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, NoUser(id)
		}
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]*User, error) {
	rows, err := s.conn.Query(ctx, getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Number, &user.Balance); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) Charge(ctx context.Context, charge *ChargeRequest) (float64, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0.00, err
	}

	defer func() { wrapErr(ctx, tx, err) }()

	var balance float64
	if err = s.conn.QueryRow(ctx, chargeQuery, charge.Amount, charge.AccountNumber).Scan(&balance); err != nil {
		return balance, err
	}

	if err = tx.Commit(ctx); err != nil {
		return balance, err
	}

	return balance, nil
}

func (s *Storage) Transfer(ctx context.Context, transfer *TransferRequest) (float64, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0.00, err
	}

	defer func() { wrapErr(ctx, tx, err) }()

	var balance float64
	if err = tx.QueryRow(ctx, balanceQuery, transfer.FromAccount).Scan(&balance); err != nil {
		return balance, err
	}

	if balance < transfer.Amount {
		return balance, InsufficientFunds(balance, transfer.Amount)
	}

	_, err = tx.Exec(ctx, transferWithdrawalQuery, transfer.Amount, transfer.FromAccount)
	if err != nil {
		return balance, err
	}

	_, err = tx.Exec(ctx, transferChargeQuery, transfer.Amount, transfer.ToAccount)
	if err != nil {
		return balance, err
	}

	if err = tx.Commit(ctx); err != nil {
		return balance, err
	}

	return balance - transfer.Amount, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	res, err := s.conn.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return NoUser(id)
	}

	return nil
}

func wrapErr(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Println(rollbackErr)
		}
	}
}

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

func (s *Storage) Register(ctx context.Context, user *User) error {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to Rollback transaction: %s\n", err)
			}
		}
	}()

	userQuery := `INSERT INTO users(first_name, last_name, email, phone_number, password_hash)
			 	  VALUES ($1, $2, $3, $4, $5)
			      RETURNING id`

	var userID int
	if err := tx.QueryRow(ctx, userQuery, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.PasswordHash).Scan(&userID); err != nil {
		return err
	}

	accountNumber := generateAccountNumber()
	accountQuery := `INSERT INTO accounts(user_id, number)
					 VALUES ($1, $2)`

	_, err = tx.Exec(ctx, accountQuery, userID, accountNumber)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*User, error) {
	query := `SELECT u.id, u.first_name, u.last_name, u.email, u.phone_number, u.password_hash, u.created_at, a.number, a.balance
			  FROM users AS u
			  JOIN accounts AS a
		  	  ON u.id = a.user_id
			  WHERE id=$1`

	rows, err := s.conn.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %d", id)
	}
	defer rows.Close()

	user := new(User)

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Number, &user.Balance); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]*User, error) {
	query := `SELECT u.id, u.first_name, u.last_name, u.email, u.phone_number, u.password_hash, u.created_at, a.number, a.balance
			  FROM users AS u
			  JOIN accounts AS a 
			  ON u.id = a.user_id`

	rows, err := s.conn.Query(ctx, query)
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

func (s *Storage) Charge(ctx context.Context, charge *ChargeRequest) error {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to Rollback transaction: %s\n", err)
			}
		}
	}()

	query := `UPDATE accounts AS a
			  SET balance = a.balance + $1
			  FROM users AS u
			  WHERE a.user_id = u.id AND a.number=$2`

	_, err = tx.Exec(ctx, query, charge.Amount, charge.AccountNumber)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) Transfer(ctx context.Context, transfer *TransferRequest) error {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to Rollback transaction: %s\n", err)
			}
		}
	}()

	balanceQuery := `SELECT a.balance
					 FROM accounts AS a
					 JOIN users AS u
					 ON a.user_id = u.id
					 WHERE a.number=$1`

	var balance float64
	if err := tx.QueryRow(ctx, balanceQuery, transfer.FromAccount).Scan(&balance); err != nil {
		return err
	}

	if balance < transfer.Amount {
		return fmt.Errorf("insufficient funds: balance %.2f, amount %.2f", balance, transfer.Amount)
	}

	withdrawalQuery := `UPDATE accounts AS a
						SET balance = a.balance - $1
						FROM users AS u
						WHERE a.user_id = u.id AND a.number=$2`

	_, err = tx.Exec(ctx, withdrawalQuery, transfer.Amount, transfer.FromAccount)
	if err != nil {
		return err
	}

	chargeQuery := `UPDATE accounts AS a
					SET balance = a.balance + $1
					FROM users AS u
					WHERE a.user_id = u.id AND a.number=$2`

	_, err = tx.Exec(ctx, chargeQuery, transfer.Amount, transfer.ToAccount)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

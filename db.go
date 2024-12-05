package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	transferTransaction = "transfer"
	chargeTransaction   = "charge"
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

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Register(ctx context.Context, u *User) (user User, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return user, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	if err = tx.QueryRow(ctx, insertNewUserQuery, u.FirstName, u.LastName, u.Email, u.PhoneNumber, u.PasswordHash, u.Acc.Number).Scan(&user.ID, &user.Acc.ID, &user.Acc.Number); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user, UserAlreadyExists(u.Email, u.PhoneNumber)
		}

		return user, err
	}

	return user, nil
}

func (s *Storage) Charge(ctx context.Context, charge *TransactionRequest) (tr Transaction, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return tr, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	_, err = tx.Exec(ctx, chargeQuery, charge.Amount, charge.ToAccount)
	if err != nil {
		return tr, err
	}

	transaction, err := insertChargeTransaction(ctx, tx, charge)
	if err != nil {
		return tr, err
	}

	return transaction, nil
}

func (s *Storage) Transfer(ctx context.Context, transfer *TransactionRequest) (tr Transaction, err error) {
	var balance float64

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return tr, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	if err = tx.QueryRow(ctx, balanceQuery, transfer.FromAccount).Scan(&balance); err != nil {
		return tr, err
	}

	if balance < transfer.Amount {
		return tr, InsufficientFunds(balance, transfer.Amount)
	}

	_, err = tx.Exec(ctx, transferQuery, transfer.FromAccount, transfer.Amount, transfer.ToAccount)
	if err != nil {
		return tr, err
	}

	transaction, err := insertTransferTransaction(ctx, tx, transfer)
	if err != nil {
		return tr, err
	}

	return transaction, nil
}

func (s *Storage) UserByID(ctx context.Context, id int) (user User, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})

	defer func() { err = rollback(ctx, tx, err) }()

	if err = s.conn.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Acc.ID, &user.Acc.Number, &user.Acc.Balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, NoUser(id)
		}
		return user, err
	}

	return user, nil
}

func (s *Storage) TransactionsByUser(ctx context.Context, id int) (trs []Transaction, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	transactions, err := getTransactions(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	res, err := s.conn.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return NoUser(id)
	}

	return nil
}

func (s *Storage) Users(ctx context.Context) ([]User, error) {
	var users []User

	rows, err := s.conn.Query(ctx, getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Acc.ID, &user.Acc.Number, &user.Acc.Balance); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func insertChargeTransaction(ctx context.Context, tx pgx.Tx, tr *TransactionRequest) (Transaction, error) {
	var (
		transactionID int
		createdAt     time.Time
	)

	if err := tx.QueryRow(ctx, insertChargeTransactionQuery, tr.ToAccount, tr.Type, tr.Amount, tr.ToAccount).Scan(&transactionID, &createdAt); err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{
		ID:        transactionID,
		Type:      tr.Type,
		Amount:    tr.Amount,
		ToAccount: tr.ToAccount,
		CreatedAt: createdAt,
	}

	return transaction, nil
}

func insertTransferTransaction(ctx context.Context, tx pgx.Tx, tr *TransactionRequest) (Transaction, error) {
	var (
		transactionID int
		createdAt     time.Time
	)

	if err := tx.QueryRow(ctx, insertTransferTransactionQuery, tr.FromAccount, tr.ToAccount, tr.Type, tr.Amount).Scan(&transactionID, &createdAt); err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{
		ID:        transactionID,
		Type:      tr.Type,
		Amount:    tr.Amount,
		ToAccount: tr.ToAccount,
		CreatedAt: createdAt,
	}

	return transaction, nil

}

func getTransactions(ctx context.Context, tx pgx.Tx, id int) ([]Transaction, error) {
	var transactions []Transaction

	rows, err := tx.Query(ctx, getTransactionsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fromAccount pgtype.Text

		transaction := Transaction{}
		if err := rows.Scan(&transaction.ID, &transaction.Type, &transaction.Amount, &fromAccount, &transaction.ToAccount, &transaction.CreatedAt); err != nil {
			return nil, err
		}

		if fromAccount.Valid {
			transaction.FromAccount = fromAccount.String
		} else {
			transaction.FromAccount = ""
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func rollback(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("rollback error: %s", rollbackErr)
		}
		return err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return fmt.Errorf("err: %s, commit error: %s", err, commitErr)
	}

	return nil
}

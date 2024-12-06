package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	transferTransaction        = "transfer"
	depositTransaction         = "deposit"
	insertUser                 = "insert_user"
	errDuplicateConstraintCode = "23505"
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

func (s *Storage) Register(ctx context.Context, user *User) (err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	stmt, err := tx.Prepare(ctx, insertUser, insertUserQuery)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, stmt.Name, user.FirstName, user.LastName, user.PhoneNumber, user.PasswordHash, user.Account.Balance, user.Account.Card.Number, user.Account.Card.CVV, user.Account.Card.ExpireTime)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == errDuplicateConstraintCode {
			return PhoneNumberAlreadyExists(user.PhoneNumber)
		}

		return err
	}

	return nil
}

func (s *Storage) Deposit(ctx context.Context, deposit *TransactionRequest) (transaction Transaction, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return transaction, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	var accountNumber string
	if err = s.conn.QueryRow(ctx, accountNumberQuery, deposit.ToCardNumber).Scan(&accountNumber); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return transaction, NoAccount(accountNumber)
		}
		return transaction, err
	}

	stmt, err := tx.Prepare(ctx, depositTransaction, depositQuery)
	if err != nil {
		return transaction, err
	}

	_, err = tx.Exec(ctx, stmt.Name, deposit.Amount, deposit.ToCardNumber)
	if err != nil {
		return transaction, err
	}

	transaction, err = insertDepositTransaction(ctx, tx, deposit)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *Storage) Transfer(ctx context.Context, transfer *TransactionRequest) (transaction Transaction, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return transaction, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	var balance float64
	if err = tx.QueryRow(ctx, balanceQuery, transfer.FromCardNumber).Scan(&balance); err != nil {
		return transaction, err
	}

	if balance < transfer.Amount {
		return transaction, InsufficientFunds(balance, transfer.Amount)
	}

	stmt, err := tx.Prepare(ctx, transferTransaction, transferQuery)
	if err != nil {
		return transaction, err
	}

	_, err = tx.Exec(ctx, stmt.Name, transfer.FromCardNumber, transfer.Amount, transfer.ToCardNumber)
	if err != nil {
		return transaction, err
	}

	transaction, err = insertTransferTransaction(ctx, tx, transfer)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *Storage) UserByID(ctx context.Context, id int) (user User, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})

	defer func() { err = rollback(ctx, tx, err) }()

	if err = s.conn.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.CreatedAt, &user.Account.ID, &user.Account.Balance, &user.Account.Card.ID, &user.Account.Card.Number, &user.Account.Card.ExpireTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, NoUser(id)
		}
		return user, err
	}

	return user, nil
}

func (s *Storage) TransactionsByUser(ctx context.Context, id int) (transactions []Transaction, err error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}

	defer func() { err = rollback(ctx, tx, err) }()

	rows, err := tx.Query(ctx, getTransactionsByUserQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		transaction := Transaction{}
		if err := rows.Scan(&transaction.ID, &transaction.Type, &transaction.Amount, &transaction.ToCardNumber, &transaction.FromCardNumber, &transaction.CreatedAt); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
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
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.CreatedAt, &user.Account.ID, &user.Account.Balance, &user.Account.Card.ID, &user.Account.Card.Number, &user.Account.Card.ExpireTime); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func insertDepositTransaction(ctx context.Context, tx pgx.Tx, tr *TransactionRequest) (Transaction, error) {
	var (
		transactionID uuid.UUID
		createdAt     time.Time
	)

	if err := tx.QueryRow(ctx, insertDepositTransactionQuery, tr.ToCardNumber, tr.Type, tr.Amount, tr.ToCardNumber).Scan(&transactionID, &createdAt); err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{
		ID:           transactionID,
		Type:         tr.Type,
		Amount:       tr.Amount,
		ToCardNumber: tr.ToCardNumber,
		CreatedAt:    createdAt,
	}

	return transaction, nil
}

func insertTransferTransaction(ctx context.Context, tx pgx.Tx, tr *TransactionRequest) (Transaction, error) {
	var (
		transactionID uuid.UUID
		createdAt     time.Time
	)

	if err := tx.QueryRow(ctx, insertTransferTransactionQuery, tr.ToCardNumber, tr.FromCardNumber, tr.Type, tr.Amount).Scan(&transactionID, &createdAt); err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{
		ID:           transactionID,
		Type:         tr.Type,
		Amount:       tr.Amount,
		ToCardNumber: tr.ToCardNumber,
		CreatedAt:    createdAt,
	}

	return transaction, nil

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

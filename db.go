package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
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

func (s *Storage) Register(ctx context.Context, user *User) (int, error) {
	var userID int

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return 0, err
	}

	defer func() { rollback(ctx, tx, err) }()

	if err = tx.QueryRow(ctx, createUserQuery, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.PasswordHash).Scan(&userID); err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, createAccountQuery, userID, user.Acc.Number)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Storage) Charge(ctx context.Context, charge *TransactionRequest) (*Transaction, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, err
	}

	defer func() { rollback(ctx, tx, err) }()

	_, err = tx.Exec(ctx, chargeQuery, charge.Amount, charge.ToAccount)
	if err != nil {
		return nil, err
	}

	transaction, err := insertTransaction(ctx, tx, charge)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Storage) Transfer(ctx context.Context, transfer *TransactionRequest) (*Transaction, error) {
	var balance float64

	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, err
	}

	defer func() { rollback(ctx, tx, err) }()

	if err = tx.QueryRow(ctx, balanceQuery, transfer.FromAccount).Scan(&balance); err != nil {
		return nil, err
	}

	if balance < transfer.Amount {
		return nil, InsufficientFunds(balance, transfer.Amount)
	}

	_, err = tx.Exec(ctx, transferWithdrawalQuery, transfer.Amount, transfer.FromAccount)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, transferChargeQuery, transfer.Amount, transfer.ToAccount)
	if err != nil {
		return nil, err
	}

	transaction, err := insertTransaction(ctx, tx, transfer)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*User, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})

	defer func() { rollback(ctx, tx, err) }()

	user := new(User)

	err = s.conn.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Acc.ID, &user.Acc.Number, &user.Acc.Balance)
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

func (s *Storage) GetTransactionsByUser(ctx context.Context, id int) ([]*Transaction, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}

	defer func() { rollback(ctx, tx, err) }()

	transactions, err := getTransactions(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
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

func (s *Storage) GetUsers(ctx context.Context) ([]*User, error) {
	var users []*User

	rows, err := s.conn.Query(ctx, getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.PasswordHash, &user.CreatedAt, &user.Acc.ID, &user.Acc.Number, &user.Acc.Balance); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func insertTransaction(ctx context.Context, tx pgx.Tx, tr *TransactionRequest) (*Transaction, error) {
	var (
		transactionID int
		accountID     int
		createdAt     time.Time
	)

	if err := tx.QueryRow(ctx, getAccountIDQuery, tr.ToAccount).Scan(&accountID); err != nil {
		return nil, err
	}

	if err := tx.QueryRow(ctx, insertTransactionQuery, accountID, tr.Type, tr.Amount, tr.FromAccount, tr.ToAccount).Scan(&transactionID, &createdAt); err != nil {
		return nil, err
	}

	transaction := &Transaction{
		ID:          transactionID,
		AccountID:   accountID,
		Type:        tr.Type,
		Amount:      tr.Amount,
		FromAccount: tr.FromAccount,
		ToAccount:   tr.ToAccount,
		CreatedAt:   createdAt,
	}

	return transaction, nil
}

func getTransactions(ctx context.Context, tx pgx.Tx, id int) ([]*Transaction, error) {
	var transactions []*Transaction

	rows, err := tx.Query(ctx, getTransactionsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fromAccount pgtype.Text

		transaction := new(Transaction)
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

func rollback(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Println(rollbackErr)
		}
		log.Println("tx rolled back")
	}
}

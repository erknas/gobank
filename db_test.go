package main

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	minAmount = 0.01
	maxAmount = 1000.00
)

type Suite struct {
	*testing.T
	store Storer
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	if err := godotenv.Load(); err != nil {
		t.Fatal(err)
	}

	t.Helper()
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	connStr := os.Getenv("DATABASE_URL")

	store, err := NewStorage(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}

	return ctx, &Suite{
		T:     t,
		store: store,
	}
}

func TestRegister(t *testing.T) {
	ctx, st := NewSuite(t)

	user := fakeUser()

	id, err := st.store.Register(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestRegister_Duplicate(t *testing.T) {
	ctx, st := NewSuite(t)

	user := fakeUser()

	id, err := st.store.Register(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	card := gofakeit.CreditCard()

	tests := []struct {
		name        string
		user        User
		expectedErr string
	}{
		{
			name: "Duplicate PhoneNumber",
			user: User{
				FirstName:    gofakeit.FirstName(),
				LastName:     gofakeit.LastName(),
				PhoneNumber:  user.PhoneNumber,
				PasswordHash: randomFakePassword(),
				Account: Account{
					Balance: user.Account.Balance,
					Card: Card{
						Number:     strconv.Itoa(card.Number),
						CVV:        card.Cvv,
						ExpireTime: card.Exp,
					},
				},
			},
			expectedErr: "user already exists",
		},
		{
			name: "Duplicate CardNumber",
			user: User{
				FirstName:    gofakeit.FirstName(),
				LastName:     gofakeit.LastName(),
				PhoneNumber:  gofakeit.Phone(),
				PasswordHash: randomFakePassword(),
				Account: Account{
					Balance: user.Account.Balance,
					Card: Card{
						Number:     user.Account.Card.Number,
						CVV:        card.Cvv,
						ExpireTime: card.Exp,
					},
				},
			},
			expectedErr: "user already exists",
		},
		{
			name: "Duplicate PhoneNumber and CardNumber",
			user: User{
				FirstName:    gofakeit.FirstName(),
				LastName:     gofakeit.LastName(),
				PhoneNumber:  user.PhoneNumber,
				PasswordHash: randomFakePassword(),
				Account: Account{
					Balance: user.Account.Balance,
					Card: Card{
						Number:     user.Account.Card.Number,
						CVV:        card.Cvv,
						ExpireTime: card.Exp,
					},
				},
			},
			expectedErr: "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := st.store.Register(ctx, &tt.user)
			require.Error(t, err)
			assert.Empty(t, id)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestDeposit(t *testing.T) {
	ctx, st := NewSuite(t)

	user := fakeUser()

	id, err := st.store.Register(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	deposit := TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: user.Account.Card.Number,
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	assert.GreaterOrEqual(t, deposit.Amount, minAmount)
	assert.LessOrEqual(t, deposit.Amount, maxAmount)

	tr, err := st.store.Deposit(ctx, &deposit)
	require.NoError(t, err)
	assert.NotEmpty(t, tr)

	u, err := st.store.UserByID(ctx, id)
	require.NoError(t, err)
	assert.NotEmpty(t, u)

	assert.Equal(t, user.Account.Balance+deposit.Amount, user.Account.Balance)
	assert.Equal(t, user.Account.Balance, user.Account.Balance-tr.Amount)

	card := gofakeit.CreditCard()
	deposit = TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: strconv.Itoa(card.Number),
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	tr, err = st.store.Deposit(ctx, &deposit)
	require.Error(t, err)
	assert.ErrorContains(t, err, "account doesn't exist")
	assert.Empty(t, tr)
}

func TestTransfer(t *testing.T) {
	ctx, st := NewSuite(t)

	user1 := fakeUser()
	user2 := fakeUser()

	id1, err := st.store.Register(ctx, user1)
	require.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := st.store.Register(ctx, user2)
	require.NoError(t, err)
	assert.NotEmpty(t, id2)

	u1, err := st.store.UserByID(ctx, id1)
	require.NoError(t, err)
	assert.NotEmpty(t, u1)

	deposit := TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: u1.Account.Card.Number,
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	depositToUser1, err := st.store.Deposit(ctx, &deposit)
	require.NoError(t, err)
	assert.NotEmpty(t, depositToUser1)

	u1, err = st.store.UserByID(ctx, id1)
	require.NoError(t, err)
	assert.NotEmpty(t, u1)

	u2, err := st.store.UserByID(ctx, id2)
	require.NoError(t, err)
	assert.NotEmpty(t, u2)

	transfer := TransactionRequest{
		Type:           transferTransaction,
		FromCardNumber: u1.Account.Card.Number,
		ToCardNumber:   u2.Account.Card.Number,
		Amount:         u1.Account.Balance,
	}

	tr, err := st.store.Transfer(ctx, &transfer)
	require.NoError(t, err)
	assert.NotEmpty(t, tr)

	u1, err = st.store.UserByID(ctx, id1)
	require.NoError(t, err)
	assert.NotEmpty(t, u1)

	u2, err = st.store.UserByID(ctx, id2)
	require.NoError(t, err)
	assert.NotEmpty(t, u2)

	assert.Empty(t, u1.Account.Balance)
	assert.Equal(t, u1.Account.Balance+transfer.Amount, u2.Account.Balance)
	assert.Equal(t, tr.FromCardNumber, u1.Account.Card.Number)
	assert.Equal(t, tr.ToCardNumber, u2.Account.Card.Number)
}

func TestUserByID(t *testing.T) {
	ctx, st := NewSuite(t)

	newUser := fakeUser()

	id, err := st.store.Register(ctx, newUser)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	user, err := st.store.UserByID(ctx, id)
	require.NoError(t, err)
	assert.NotEmpty(t, user)

	fakeId := gofakeit.Uint8()

	fakeUser, err := st.store.UserByID(ctx, int(fakeId))
	require.Error(t, err)
	assert.Empty(t, fakeUser)
	assert.ErrorContains(t, err, "user doesn't exist")
}

func TestTransactionsByUser(t *testing.T) {
	ctx, st := NewSuite(t)

	newUser1 := fakeUser()
	newUser2 := fakeUser()

	id1, err := st.store.Register(ctx, newUser1)
	require.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := st.store.Register(ctx, newUser2)
	require.NoError(t, err)
	assert.NotEmpty(t, id1)

	user1, err := st.store.UserByID(ctx, id1)
	require.NoError(t, err)
	assert.NotEmpty(t, user1)

	user2, err := st.store.UserByID(ctx, id2)
	require.NoError(t, err)
	assert.NotEmpty(t, user2)

	deposit1 := TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: user1.Account.Card.Number,
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	tr1, err := st.store.Deposit(ctx, &deposit1)
	require.NoError(t, err)
	assert.NotEmpty(t, tr1)

	deposit2 := TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: user1.Account.Card.Number,
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	tr2, err := st.store.Deposit(ctx, &deposit2)
	require.NoError(t, err)
	assert.NotEmpty(t, tr2)

	transfer := TransactionRequest{
		Type:           transferTransaction,
		FromCardNumber: user1.Account.Card.Number,
		ToCardNumber:   user2.Account.Card.Number,
		Amount:         1,
	}

	tr3, err := st.store.Transfer(ctx, &transfer)
	require.NoError(t, err)
	assert.NotEmpty(t, tr3)

	trs, err := st.store.TransactionsByUser(ctx, user1.ID)
	require.NoError(t, err)
	assert.Len(t, trs, 3)
}

func TestInsertDepositTransaction(t *testing.T) {
	ctx, st := NewSuite(t)

	newUser := fakeUser()

	id, err := st.store.Register(ctx, newUser)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	user, err := st.store.UserByID(ctx, id)
	require.NoError(t, err)
	assert.NotEmpty(t, user)

	deposit := TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: user.Account.Card.Number,
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	tr, err := st.store.Deposit(ctx, &deposit)
	require.NoError(t, err)
	assert.NotEmpty(t, tr)

	now := time.Now().UTC()
	delta := time.Second

	assert.NotEmpty(t, tr.ID)
	assert.Empty(t, tr.FromCardNumber)
	assert.Equal(t, deposit.Amount, tr.Amount)
	assert.Greater(t, tr.Amount, 0.00)
	assert.Equal(t, user.Account.Card.Number, tr.ToCardNumber)
	assert.WithinDuration(t, now, tr.CreatedAt, delta)
}

func TestInsertTransferTransaction(t *testing.T) {
	ctx, st := NewSuite(t)

	newUser1 := fakeUser()
	newUser2 := fakeUser()

	id1, err := st.store.Register(ctx, newUser1)
	require.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := st.store.Register(ctx, newUser2)
	require.NoError(t, err)
	assert.NotEmpty(t, id2)

	user1, err := st.store.UserByID(ctx, id1)
	require.NoError(t, err)
	assert.NotEmpty(t, user1)

	user2, err := st.store.UserByID(ctx, id2)
	require.NoError(t, err)
	assert.NotEmpty(t, user2)

	deposit := TransactionRequest{
		Type:         depositTransaction,
		ToCardNumber: user1.Account.Card.Number,
		Amount:       gofakeit.Price(minAmount, maxAmount),
	}

	depositTr, err := st.store.Deposit(ctx, &deposit)
	require.NoError(t, err)
	assert.NotEmpty(t, depositTr)

	user1, err = st.store.UserByID(ctx, id1)
	require.NoError(t, err)
	assert.NotEmpty(t, user2)

	transfer := TransactionRequest{
		Type:           transferTransaction,
		FromCardNumber: user1.Account.Card.Number,
		ToCardNumber:   user2.Account.Card.Number,
		Amount:         user1.Account.Balance,
	}

	tr, err := st.store.Transfer(ctx, &transfer)
	require.NoError(t, err)
	assert.NotEmpty(t, tr)

	user1Tr, err := st.store.TransactionsByUser(ctx, user1.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, user1Tr)

	user2Tr, err := st.store.TransactionsByUser(ctx, user2.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, user2Tr)

	now := time.Now().UTC()
	delta := time.Second

	assert.WithinDuration(t, now, user1Tr[0].CreatedAt, delta)
	assert.WithinDuration(t, now, user2Tr[0].CreatedAt, delta)
	assert.Equal(t, user1Tr[0].CreatedAt, user2Tr[0].CreatedAt)
	assert.Equal(t, user1Tr[0].ID, user2Tr[0].ID)
	assert.Equal(t, user1Tr[0].FromCardNumber, transfer.FromCardNumber)
	assert.Equal(t, user2Tr[0].ToCardNumber, transfer.ToCardNumber)
	assert.Equal(t, user1Tr[0].Type, transferTransaction)
	assert.Equal(t, user2Tr[0].Type, transferTransaction)
	assert.Equal(t, transfer.Amount, user1Tr[0].Amount)
	assert.Equal(t, transfer.Amount, user2Tr[0].Amount)
}

func fakeUser() *User {
	firstName := gofakeit.Name()
	lastName := gofakeit.LastName()
	phoneNumber := gofakeit.Phone()
	password := randomFakePassword()

	user, err := NewUser(&NewUserRequest{FirstName: firstName, LastName: lastName, PhoneNumber: phoneNumber, Password: password})
	if err != nil {
		return nil
	}

	return user
}

package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const (
	cardNumberLength  = 16
	cvvLength         = 3
	phoneNumberLength = 10
	initialBalance    = 0.00
)

func TestNewCard(t *testing.T) {
	card := NewCard()

	assert.Len(t, card.Number, cardNumberLength)
	assert.Len(t, card.CVV, cvvLength)
	assert.Regexp(t, regexp.MustCompile(`^[0-9]+$`), card.Number)
	assert.Regexp(t, regexp.MustCompile(`^[0-9]+$`), card.CVV)

	expireTime := time.Now().AddDate(0, 61, 0).Format(layout)
	assert.Equal(t, expireTime, card.ExpireTime)
}

func TestNewUser(t *testing.T) {
	firstName := gofakeit.Name()
	lastName := gofakeit.LastName()
	phoneNumber := gofakeit.Phone()
	password := randomFakePassword()

	user, err := NewUser(&NewUserRequest{FirstName: firstName, LastName: lastName, PhoneNumber: phoneNumber, Password: password})
	require.NoError(t, err)

	now := time.Now().UTC()
	delta := time.Millisecond

	assert.NotEmpty(t, user.FirstName)
	assert.NotEmpty(t, user.LastName)
	assert.Len(t, user.PhoneNumber, phoneNumberLength)
	assert.Regexp(t, regexp.MustCompile(`^[0-9]+$`), phoneNumber)
	assert.NotEmpty(t, user.PasswordHash)
	assert.WithinDuration(t, now, user.CreatedAt, delta)
	assert.Equal(t, initialBalance, user.Account.Balance)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	require.NoError(t, err)
}

func TestNewUser_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		user        NewUserRequest
		expectedErr string
	}{
		{
			name: "Empty FirstName",
			user: NewUserRequest{
				FirstName:   "",
				LastName:    gofakeit.LastName(),
				PhoneNumber: gofakeit.Phone(),
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Empty LastName",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    "",
				PhoneNumber: gofakeit.Phone(),
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Empty PhoneNumber",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				PhoneNumber: "",
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Empty Password",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				PhoneNumber: gofakeit.Phone(),
				Password:    "",
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Empty all",
			user: NewUserRequest{
				FirstName:   "",
				LastName:    "",
				PhoneNumber: "",
				Password:    "",
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Invalid PhoneNumber length",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				PhoneNumber: "897651",
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Invalid PhoneNumber contains letters",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				PhoneNumber: "m89u7699a1",
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Invalid PhoneNumber contains symbols",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				PhoneNumber: "8@#$%^&*4)",
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
		{
			name: "Invalid PhoneNumber contains letters, symbols, digits",
			user: NewUserRequest{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				PhoneNumber: "910-a4 j/3",
				Password:    randomFakePassword(),
			},
			expectedErr: "invalid user data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(&tt.user)
			require.Error(t, err)
			require.Empty(t, user)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 10)
}

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
	password := gofakeit.Password(true, true, true, false, false, 10)

	user, err := NewUser(firstName, lastName, phoneNumber, password)
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

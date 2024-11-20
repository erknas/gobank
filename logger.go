package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type RequestID struct{}

type Logger struct {
	next Storer
}

func NewLogger(next Storer) *Logger {
	return &Logger{
		next: next,
	}
}

func (l *Logger) Register(ctx context.Context, user *User) (err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"time":       begin,
			"took":       time.Since(begin),
			"request_id": ctx.Value(RequestID{}),
			"error":      err,
		}).Info("register user")
	}(time.Now())

	return l.next.Register(ctx, user)
}

func (l *Logger) GetUserByID(ctx context.Context, id int) (user *User, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"time":       begin,
			"took":       time.Since(begin),
			"request_id": ctx.Value(RequestID{}),
			"error":      err,
			"user":       fmt.Sprintf("%+v", user),
		}).Info("get user")
	}(time.Now())

	return l.next.GetUserByID(ctx, id)
}

func (l *Logger) GetUsers(ctx context.Context) (users []*User, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"time":       begin,
			"took":       time.Since(begin),
			"request_id": ctx.Value(RequestID{}),
			"error":      err,
			"users":      len(users),
		}).Info("get users")
	}(time.Now())

	return l.next.GetUsers(ctx)
}

func (l *Logger) Charge(ctx context.Context, charge *ChargeRequest) (err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"time":           begin,
			"took":           time.Since(begin),
			"request_id":     ctx.Value(RequestID{}),
			"error":          err,
			"account number": charge.AccountNumber,
			"amount":         charge.Amount,
		}).Info("charge")
	}(time.Now())

	return l.next.Charge(ctx, charge)
}

func (l *Logger) Transfer(ctx context.Context, transfer *TransferRequest) (err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"time":       begin,
			"took":       time.Since(begin),
			"request_id": ctx.Value(RequestID{}),
			"error":      err,
		}).Info("transfer")
	}(time.Now())

	return l.next.Transfer(ctx, transfer)
}

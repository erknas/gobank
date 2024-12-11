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
	log  *logrus.Logger
}

func NewLogger(next Storer) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.DebugLevel)

	return &Logger{
		next: next,
		log:  log,
	}
}

func (l *Logger) Register(ctx context.Context, user *User) (id int, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("register user")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
			}).Error("register user failed")
		}
	}(time.Now())

	return l.next.Register(ctx, user)
}

func (l *Logger) Deposit(ctx context.Context, deposit *TransactionRequest) (transaction Transaction, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("deposit")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
			}).Error("deposit failed")
		}
	}(time.Now())

	return l.next.Deposit(ctx, deposit)
}

func (l *Logger) Transfer(ctx context.Context, transfer *TransactionRequest) (transaction Transaction, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("transfer")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
			}).Error("transfer failed")
		}
	}(time.Now())

	return l.next.Transfer(ctx, transfer)
}

func (l *Logger) UserByID(ctx context.Context, id int) (user User, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("get user")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
				"user ID":    id,
			}).Error("get user failed")
		}
	}(time.Now())

	return l.next.UserByID(ctx, id)
}

func (l *Logger) TransactionsByUser(ctx context.Context, id int) (transactions []Transaction, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("get transactions by user")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
				"user ID":    id,
			}).Error("get transactions by user failed")
		}
	}(time.Now())

	return l.next.TransactionsByUser(ctx, id)
}

func (l *Logger) Users(ctx context.Context) (users []User, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("get users")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
			}).Error("get users failed")
		}
	}(time.Now())

	return l.next.Users(ctx)
}

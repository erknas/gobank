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
				"user ID":    user.ID,
			}).Error("register user failed")
		}
	}(time.Now())

	return l.next.Register(ctx, user)
}

func (l *Logger) Charge(ctx context.Context, charge *TransactionRequest) (transaction *Transaction, err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("charge")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
			}).Error("charge failed")
		}
	}(time.Now())

	return l.next.Charge(ctx, charge)
}

func (l *Logger) Transfer(ctx context.Context, transfer *TransactionRequest) (transaction *Transaction, err error) {
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

func (l *Logger) GetUserByID(ctx context.Context, id int) (user *User, err error) {
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

	return l.next.GetUserByID(ctx, id)
}

func (l *Logger) GetTransactionsByUser(ctx context.Context, id int) (transactions []*Transaction, err error) {
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

	return l.next.GetTransactionsByUser(ctx, id)
}

func (l *Logger) Delete(ctx context.Context, id int) (err error) {
	defer func(begin time.Time) {
		if err == nil {
			l.log.WithFields(logrus.Fields{
				"took":       fmt.Sprintf("%dµs", time.Since(begin).Microseconds()),
				"request_id": ctx.Value(RequestID{}),
			}).Info("delete")
		} else {
			l.log.WithFields(logrus.Fields{
				"request_id": ctx.Value(RequestID{}),
				"error":      err,
				"user ID":    id,
			}).Error("delete failed")
		}
	}(time.Now())

	return l.next.Delete(ctx, id)
}

func (l *Logger) GetUsers(ctx context.Context) (users []*User, err error) {
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

	return l.next.GetUsers(ctx)
}

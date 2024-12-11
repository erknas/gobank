package main

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid JSON request"))
}

func InvalidID() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid user ID"))
}

func NoUser() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("user doesn't exist"))
}

func NoAccount() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("account doesn't exist"))
}

func UserExists() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("user already exists"))
}

func InsufficientFunds(balance, amount float64) APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("insufficient funds: balance = %.2f, amount = %.2f", balance, amount))
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

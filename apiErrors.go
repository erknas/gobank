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
	return fmt.Sprintf("api error: %d", e.StatusCode)
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

func NoUser(id int) APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("user doesn't exist: ID = %d", id))
}

func NoAccount(accountNumber string) APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("account doesn't exist with number %s", accountNumber))
}

func PhoneNumberAlreadyExists(phoneNumber string) APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("phone number %s already exists", phoneNumber))
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

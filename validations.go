package main

import (
	"fmt"
	"unicode"
)

func (r NewUserRequest) ValidateUserData() map[string]string {
	errors := make(map[string]string)

	if len(r.FirstName) == 0 {
		errors["firstName"] = "fist name should not be empty"
	}

	if len(r.LastName) == 0 {
		errors["lastName"] = "last name should not be empty"
	}

	if len(r.PhoneNumber) != 10 {
		errors["phoneNumber"] = "invalid phone number"
	}

	if len(r.Password) == 0 {
		errors["password"] = "password should not be empty"
	}

	return errors
}

func (r TransactionRequest) ValidateTransaction() map[string]string {
	errors := make(map[string]string)

	if r.Type != transferTransaction && r.Type != depositTransaction {
		errors["transaction type"] = "unsupported transaction"
	}

	if r.Type == transferTransaction {

		if len(r.FromCardNumber) != 16 {
			errors["fromCardNumber"] = fmt.Sprintf("invalid card number: length should be 16, got %d", len(r.FromCardNumber))
		}

		if len(r.ToCardNumber) != 16 {
			errors["toCardNumber"] = fmt.Sprintf("invalid card number: length should be 16, got %d", len(r.ToCardNumber))
		}

		for _, digit := range r.FromCardNumber {
			if !unicode.IsDigit(digit) {
				errors["fromCardNumber"] = "card number should contains only digits"
				break
			}
		}

		for _, digit := range r.ToCardNumber {
			if !unicode.IsDigit(digit) {
				errors["toCardNumber"] = "card number should contains only digits"
				break
			}
		}

		if r.Amount < 0 {
			errors["amount"] = "amount cannot be negative"
		}

		if r.Amount == 0 {
			errors["amount"] = "amount cannot be zero"
		}
	}

	if r.Type == depositTransaction {

		if len(r.ToCardNumber) != 16 {
			errors["accountNumber"] = fmt.Sprintf("invalid card number: length should be 16, got %d", len(r.ToCardNumber))
		}

		for _, digit := range r.ToCardNumber {
			if !unicode.IsDigit(digit) {
				errors["toCardNumber"] = "card number should contains only digits"
				break
			}
		}

		if r.Amount < 0 {
			errors["amount"] = "amount cannot be negative"
		}

		if r.Amount == 0 {
			errors["amount"] = "amount cannot be zero"
		}
	}

	return errors
}

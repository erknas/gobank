package main

func (r NewUserRequest) ValidateUserData() map[string]string {
	errors := make(map[string]string)

	if len(r.FirstName) == 0 {
		errors["firstName"] = "fist name should not be empty"
	}

	if len(r.LastName) == 0 {
		errors["lastName"] = "last name should not be empty"
	}

	if len(r.PhoneNumber) != 11 {
		errors["phoneNumber"] = "wrong phone number"
	}

	if len(r.Password) == 0 {
		errors["password"] = "password should not be empty"
	}

	return errors
}

func (r TransactionRequest) ValidateTransaction() map[string]string {
	errors := make(map[string]string)

	if r.Type != "transfer" && r.Type != "charge" {
		errors["transaction type"] = "unsupported transaction"
	}

	if r.Type == "transfer" {
		if r.Amount < 0 {
			errors["amount"] = "amount cannot be negative"
		}

		if r.Amount == 0 {
			errors["amount"] = "amount cannot be zero"
		}

		if r.FromAccount == "" {
			errors["fromAccount"] = "provide account number"
		}

		if r.ToAccount == "" {
			errors["toAccount"] = "provide account number"
		}
	}

	if r.Type == "charge" {
		if r.Amount < 0 {
			errors["amount"] = "amount cannot be negative"
		}

		if r.Amount == 0 {
			errors["amount"] = "amount cannot be zero"
		}

		if r.ToAccount == "" {
			errors["toAccount"] = "provide account number"
		}
	}

	return errors
}

package main

func (r RegisterUserRequest) ValidateRequestData() map[string]string {
	errors := make(map[string]string, 5)

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

package main

var (
	createUserQuery         = `INSERT INTO users(first_name, last_name, email, phone_number, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	createAccountQuery      = `INSERT INTO accounts(user_id, number) VALUES ($1, $2)`
	getUserByIDQuery        = `SELECT u.id, u.first_name, u.last_name, u.email, u.phone_number, u.password_hash, u.created_at, a.number, a.balance FROM users AS u JOIN accounts AS a ON u.id = a.user_id WHERE u.id = $1`
	getUsersQuery           = `SELECT u.id, u.first_name, u.last_name, u.email, u.phone_number, u.password_hash, u.created_at, a.number, a.balance FROM users AS u JOIN accounts AS a ON u.id = a.user_id`
	chargeQuery             = `UPDATE accounts AS a SET balance = a.balance + $1 FROM users AS u WHERE a.user_id = u.id AND a.number = $2 RETURNING balance`
	balanceQuery            = `SELECT a.balance FROM accounts AS a JOIN users AS u ON a.user_id = u.id WHERE a.number = $1`
	transferWithdrawalQuery = `UPDATE accounts AS a SET balance = a.balance - $1 FROM users AS u WHERE a.user_id = u.id AND a.number = $2`
	transferChargeQuery     = `UPDATE accounts AS a SET balance = a.balance + $1 FROM users AS u WHERE a.user_id = u.id AND a.number = $2`
	deleteUserQuery         = `DELETE FROM users AS u WHERE u.id = $1`
	getAccountIDQuery       = `SELECT a.id FROM accounts AS a JOIN users AS u ON a.user_id = u.id WHERE a.number = $1`
	insertTransactionQuery  = `INSERT INTO transactions(account_id, transaction_type, amount, from_account_number, to_account_number) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	getTransactionsQuery    = `SELECT id, account_id, transaction_type, amount, from_account_number, to_account_number, created_at FROM transactions`
)

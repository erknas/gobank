package main

var (
	createUserQuery                = `INSERT INTO users(first_name, last_name, email, phone_number, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	createAccountQuery             = `INSERT INTO accounts(user_id, number) VALUES ($1, $2)`
	getUserByIDQuery               = `SELECT u.id, u.first_name, u.last_name, u.email, u.phone_number, u.password_hash, u.created_at, a.id, a.number, a.balance FROM users AS u JOIN accounts AS a ON u.id = a.user_id WHERE u.id = $1`
	getUsersQuery                  = `SELECT u.id, u.first_name, u.last_name, u.email, u.phone_number, u.password_hash, u.created_at, a.id, a.number, a.balance FROM users AS u JOIN accounts AS a ON u.id = a.user_id`
	chargeQuery                    = `UPDATE accounts AS a SET balance = a.balance + $1 FROM users AS u WHERE a.user_id = u.id AND a.number = $2 RETURNING balance`
	balanceQuery                   = `SELECT a.balance FROM accounts AS a JOIN users AS u ON a.user_id = u.id WHERE a.number = $1`
	transferQuery                  = `UPDATE accounts AS a SET balance = CASE WHEN a.number = $1 THEN a.balance - $2 WHEN a.number = $3 THEN a.balance + $2 END WHERE a.number IN ($1, $3)`
	deleteUserQuery                = `DELETE FROM users AS u WHERE u.id = $1`
	insertChargeTransactionQuery   = `WITH account AS (SELECT a.id FROM accounts AS a JOIN users AS u ON a.user_id = u.id WHERE a.number = $1) INSERT INTO transactions(account_id, transaction_type, amount, to_account_number) VALUES ((SELECT id FROM account), $2, $3, $4) RETURNING id, created_at`
	getTransactionsQuery           = `SELECT id, transaction_type, amount, from_account_number, to_account_number, created_at FROM transactions WHERE account_id = $1 ORDER BY created_at DESC`
	insertTransferTransactionQuery = `WITH from_account AS (SELECT a.id FROM accounts AS a JOIN users AS u ON a.user_id = u.id WHERE a.number = $1), to_account AS (SELECT a.id FROM accounts AS a JOIN users AS u ON a.user_id = u.id WHERE a.number = $2), inserted_from_transaction AS (INSERT INTO transactions(account_id, transaction_type, amount, from_account_number, to_account_number) VALUES ((SELECT id FROM from_account), $3, $4, $1, $2) RETURNING id, created_at) INSERT INTO transactions(account_id, transaction_type, amount, from_account_number, to_account_number) VALUES ((SELECT id FROM to_account), $3, $4, $1, $2) RETURNING id, created_at`
)

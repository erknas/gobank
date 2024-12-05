CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	email VARCHAR(100) UNIQUE,
	phone_number VARCHAR(11) UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS accounts (
	id SERIAL PRIMARY KEY,	
	user_id INT NOT NULL,
	number VARCHAR(19) NOT NULL UNIQUE,
	balance NUMERIC(15,2) DEFAULT 0.00 CHECK (balance >= 0),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);

CREATE TABLE IF NOT EXISTS transactions (
	id SERIAL PRIMARY KEY,
	account_id 	INT NOT NULL,
	transaction_type VARCHAR(10) NOT NULL,
	amount NUMERIC(9,2) NOT NULL CHECK (amount >= 0),
	to_account_number VARCHAR(19) NOT NULL,
	from_account_number VARCHAR(19),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
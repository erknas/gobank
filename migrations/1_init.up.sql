CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	phone_number VARCHAR(10) UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS accounts (
	id SERIAL PRIMARY KEY,	
	user_id INT NOT NULL,
	balance NUMERIC(15,2) DEFAULT 0.00 CHECK (balance >= 0),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);

CREATE TABLE IF NOT EXISTS cards (
	id SERIAL PRIMARY KEY,
	account_id INT NOT NULL,
	card_number VARCHAR(16) NOT NULL UNIQUE CHECK (LENGTH(card_number) = 16),
	cvv VARCHAR(3) NOT NULL CHECK (LENGTH(cvv) = 3),
	expire_time VARCHAR(5) NOT NULL,
	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID DEFAULT gen_random_uuid(),
    account_id INT NOT NULL,
    transaction_type VARCHAR(10) NOT NULL,
    amount NUMERIC(6,2) NOT NULL CHECK (amount >= 0),
    to_card_number VARCHAR(16) NOT NULL,
    from_card_number VARCHAR(16) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_transaction_id ON transactions(transaction_id);
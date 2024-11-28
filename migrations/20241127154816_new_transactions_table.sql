-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
	id SERIAL PRIMARY KEY,
	account_id 	INT NOT NULL,
	transaction_type VARCHAR(10) NOT NULL,
	amount NUMERIC(9,2) NOT NULL CHECK (amount >= 0),
	to_account_number VARCHAR(19) NOT NULL,
	from_account_number VARCHAR(19),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd

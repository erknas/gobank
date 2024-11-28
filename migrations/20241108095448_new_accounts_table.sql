-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
	id SERIAL PRIMARY KEY,	
	user_id INT NOT NULL,
	number VARCHAR(19) NOT NULL UNIQUE,
	balance NUMERIC(15,2) DEFAULT 0.00 CHECK (balance >= 0),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd

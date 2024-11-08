-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
	user_id INT NOT NULL,
	number INT NOT NULL,
	balance INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd

package main

var (
	balanceQuery = `SELECT balance 
					FROM accounts 
					JOIN users ON accounts.user_id = users.id
					JOIN cards ON accounts.id = cards.account_id
					WHERE card_number = $1;`

	accountNumberQuery = `SELECT users.id 
						  FROM users 
						  JOIN accounts ON users.id = accounts.user_id
						  JOIN cards ON accounts.id = cards.account_id 
						  WHERE cards.card_number = $1;`

	transferQuery = `UPDATE accounts 
					 SET balance = 
					 CASE 
					 	WHEN accounts.id = (SELECT cards.account_id FROM cards WHERE cards.card_number = $1) THEN balance - $2 
						WHEN accounts.id = (SELECT cards.account_id FROM cards WHERE cards.card_number = $3) THEN balance + $2 
						ELSE balance
					 END 
					 WHERE accounts.id IN (
					 (SELECT cards.account_id FROM cards WHERE cards.card_number = $1),
					 (SELECT cards.account_id FROM cards WHERE cards.card_number = $3)
					 );`

	deleteUserQuery = `DELETE 
					   FROM users 
					   WHERE users.id = $1;`

	insertTransferTransactionQuery = `WITH from_card_number_transfer AS (
									  SELECT accounts.id FROM accounts 
									  JOIN users ON accounts.user_id = users.id 
									  JOIN cards ON accounts.id = cards.account_id
	  								  WHERE cards.card_number = $1
	  								  ), 
	                                  to_card_number_transfer AS (
									  SELECT accounts.id FROM accounts 
	                                  JOIN users ON accounts.user_id = users.id 
									  JOIN cards ON accounts.id = cards.account_id
	                                  WHERE cards.card_number = $2
	                                  ), 
									  insert_transaction AS (
									  INSERT INTO transactions(account_id, transaction_type, amount, to_card_number, from_card_number) 
									  VALUES((SELECT id FROM from_card_number_transfer), $3, $4, $1, $2) RETURNING transaction_id) 
									  INSERT INTO transactions(transaction_id, account_id, transaction_type, amount, to_card_number, from_card_number) 
								      VALUES((SELECT transaction_id FROM insert_transaction), (SELECT id FROM to_card_number_transfer), $3, $4, $1, $2) RETURNING transaction_id, created_at;`

	getTransactionsByUserQuery = `SELECT transaction_id, transaction_type, amount, to_card_number, from_card_number, transactions.created_at FROM transactions
								  JOIN accounts ON transactions.account_id = accounts.id
								  JOIN users ON accounts.user_id = users.id
							      WHERE users.id = $1 
							      ORDER BY created_at DESC;`

	insertUserQuery = `WITH new_user AS (
					   INSERT INTO users(first_name, last_name, phone_number, password_hash)
					   VALUES($1, $2, $3, $4)
					   RETURNING id
					   ),
					   new_account AS (
					   INSERT INTO accounts(user_id, balance)
					   VALUES((SELECT id FROM new_user), $5)
					   RETURNING id
					   )
					   INSERT INTO cards(account_id, card_number, cvv, expire_time)
					   VALUES((SELECT id from new_account), $6, $7, $8);`

	getUserByIDQuery = `SELECT u.id, u.first_name, u.last_name, u.phone_number, u.created_at, a.id, a.balance, c.id, c.card_number, c.expire_time
						FROM users AS u
	                    JOIN accounts AS a ON u.id = a.user_id
						JOIN cards AS c ON a.id = c.account_id
						WHERE u.id = $1;`

	getUsersQuery = `SELECT users.id, users.first_name, users.last_name, users.phone_number, users.created_at, accounts.id, accounts.balance, cards.id, cards.card_number, cards.expire_time
					 FROM users
	                 JOIN accounts ON users.id = accounts.user_id
				     JOIN cards ON accounts.id = cards.account_id;`

	depositQuery = `UPDATE accounts
					SET balance = balance + $1
					WHERE id IN (
					SELECT accounts.id
					FROM accounts
					JOIN users ON accounts.user_id = users.id
					JOIN cards ON accounts.id = cards.account_id
					WHERE cards.card_number = $2);`

	insertDepositTransactionQuery = `WITH account AS (
									 SELECT accounts.id 
									 FROM accounts
									 JOIN users ON accounts.user_id = users.id
									 JOIN cards ON accounts.id = cards.account_id
									 WHERE cards.card_number = $1
									 )
									 INSERT INTO transactions(account_id, transaction_type, amount, to_card_number) 
									 VALUES((SELECT id from account LIMIT 1), $2, $3, $4) RETURNING transaction_id, created_at;`
)

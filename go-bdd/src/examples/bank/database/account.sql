CREATE TABLE accounts (
	 id SERIAL PRIMARY KEY
	,email TEXT NOT NULL
	,balance INTEGER NOT NULL DEFAULT 0
	,UNIQUE(email)
);

CREATE TABLE api_keys (
	 id SERIAL PRIMARY KEY
	,key TEXT NOT NULL
	,account_id INTEGER
	,CONSTRAINT f_api_keys_account_id FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

INSERT INTO accounts (email, balance)
VALUES
('wow@holyhell.com', 4000);

INSERT INTO api_keys (key, account_id)
VALUES
('dev', 1);

CREATE TABLE accounts (
	 id SERIAL PRIMARY KEY
	,email TEXT NOT NULL
	,balance INTEGER NOT NULL DEFAULT 0
	,UNIQUE(email)
);

INSERT INTO accounts (email, balance)
VALUES
('wow@holyhell.com', 4000);

TRUNCATE accounts RESTART IDENTITY CASCADE;

INSERT INTO accounts (email, balance)
VALUES ('wow@holyhell.com', 6000);

INSERT INTO api_keys (key, account_id)
VALUES ('dev', 1);

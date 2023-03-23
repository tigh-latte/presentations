package apikey

const (
	SelectApiKey = `
	SELECT id, key, account_id
	FROM api_keys
	WHERE key = $1
	`
)

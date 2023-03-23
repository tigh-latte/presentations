package repo

import (
	"context"
	"database/sql"

	"github.com/tigh-latte/presentations/go-bdd/src/examples/bank"
	"github.com/tigh-latte/presentations/go-bdd/src/examples/bank/repo/apikey"
)

type apikeyRepo struct {
	db *sql.DB
}

func NewAPIKeyRepo(db *sql.DB) bank.APIKeyRepo {
	return &apikeyRepo{db: db}
}

func (a *apikeyRepo) APIKey(ctx context.Context, args bank.GetAPIKeyArgs) (apiKey bank.APIKey, err error) {
	res := a.db.QueryRowContext(ctx, apikey.SelectApiKey, args.APIKey)
	err = res.Scan(&apiKey.ID, &apiKey.APIKey, &apiKey.AccountID)
	return apiKey, err
}

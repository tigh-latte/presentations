package bank

import "context"

type APIKey struct {
	ID        int64  `json:"-"`
	APIKey    string `json:"-"`
	AccountID int64  `json:"-"`
}

type GetAPIKeyArgs struct {
	APIKey string
}

type APIKeyRepo interface {
	APIKey(ctx context.Context, args GetAPIKeyArgs) (APIKey, error)
}

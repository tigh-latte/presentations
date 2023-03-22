package bank

import (
	"context"

	"gopkg.in/guregu/null.v3"
)

type Account struct {
	ID      int64  `json:"id"`
	Email   string `json:"email"`
	Balance int64  `json:"balance"`
}

type GetAccountArgs struct {
	ID int64 `db:"id"`
}

type UpdateAccountReq struct {
	ID      int64
	Email   null.String
	Balance null.Int
}

type DeleteAccountReq struct {
	ID int64
}

type AccountWithdrawArgs struct {
	ID     int64
	Amount int64
}

type AccountService interface {
	Withdraw(ctx context.Context, args AccountWithdrawArgs) (int64, error)
	AccountStore
}

type AccountStore interface {
	List(ctx context.Context) ([]Account, error)
	GetByID(ctx context.Context, args GetAccountArgs) (Account, error)
	Create(ctx context.Context, acc Account) (Account, error)
	Update(ctx context.Context, req UpdateAccountReq) (Account, error)
	Delete(ctx context.Context, req DeleteAccountReq) error
}

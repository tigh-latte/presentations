package service

import (
	"context"

	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/accountmanager"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/accountmanager/errs"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type accountSvc struct {
	repo accountmanager.AccountStore
}

func NewAccountService(repo accountmanager.AccountStore) accountmanager.AccountService {
	return &accountSvc{
		repo: repo,
	}
}

func (a *accountSvc) List(ctx context.Context) ([]accountmanager.Account, error) {
	return a.repo.List(ctx)
}

func (a *accountSvc) Create(ctx context.Context, acc accountmanager.Account) (accountmanager.Account, error) {
	return a.repo.Create(ctx, acc)
}

func (a *accountSvc) Withdraw(ctx context.Context, args accountmanager.AccountWithdrawArgs) (int64, error) {
	account, err := a.GetByID(ctx, accountmanager.GetAccountArgs{
		ID: args.ID,
	})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get account '%s'", args.ID)
	}
	if account.Balance < args.Amount {
		return 0, errors.Wrapf(errs.ErrInsufficientFunds, "error withdrawing for account '%s'", args.ID)
	}
	account.Balance -= args.Amount

	if account, err = a.Update(ctx, accountmanager.UpdateAccountReq{
		ID:      account.ID,
		Balance: null.NewInt(account.Balance, true),
	}); err != nil {
		return 0, errors.Wrapf(err, "failed to set new account balance for '%s'", account.ID)
	}

	return args.Amount, nil
}

func (a *accountSvc) GetByID(ctx context.Context, args accountmanager.GetAccountArgs) (accountmanager.Account, error) {
	return a.repo.GetByID(ctx, args)
}

func (a *accountSvc) Update(ctx context.Context, req accountmanager.UpdateAccountReq) (accountmanager.Account, error) {
	return a.repo.Update(ctx, req)
}

func (a *accountSvc) Delete(ctx context.Context, req accountmanager.DeleteAccountReq) error {
	return a.repo.Delete(ctx, req)
}

package service

import (
	"context"

	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/errs"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type accountSvc struct {
	repo bank.AccountStore
}

func NewAccountService(repo bank.AccountStore) bank.AccountService {
	return &accountSvc{
		repo: repo,
	}
}

func (a *accountSvc) List(ctx context.Context) ([]bank.Account, error) {
	return a.repo.List(ctx)
}

func (a *accountSvc) Create(ctx context.Context, acc bank.Account) (bank.Account, error) {
	return a.repo.Create(ctx, acc)
}

func (a *accountSvc) Withdraw(ctx context.Context, args bank.AccountWithdrawArgs) (int64, error) {
	account, err := a.GetByID(ctx, bank.GetAccountArgs{
		ID: args.ID,
	})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get account '%s'", args.ID)
	}
	if account.Balance < args.Amount {
		return 0, errors.Wrapf(errs.ErrInsufficientFunds, "error withdrawing for account '%s'", args.ID)
	}
	account.Balance -= args.Amount

	if account, err = a.Update(ctx, bank.UpdateAccountReq{
		ID:      account.ID,
		Balance: null.NewInt(account.Balance, true),
	}); err != nil {
		return 0, errors.Wrapf(err, "failed to set new account balance for '%s'", account.ID)
	}

	return args.Amount, nil
}

func (a *accountSvc) GetByID(ctx context.Context, args bank.GetAccountArgs) (bank.Account, error) {
	return a.repo.GetByID(ctx, args)
}

func (a *accountSvc) Update(ctx context.Context, req bank.UpdateAccountReq) (bank.Account, error) {
	return a.repo.Update(ctx, req)
}

func (a *accountSvc) Delete(ctx context.Context, req bank.DeleteAccountReq) error {
	return a.repo.Delete(ctx, req)
}

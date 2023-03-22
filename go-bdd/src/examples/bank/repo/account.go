package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/errs"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/bank/repo/account"
	"github.com/lib/pq"
)

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) bank.AccountStore {
	return &accountRepo{
		db: db,
	}
}

func (a *accountRepo) List(ctx context.Context) ([]bank.Account, error) {
	rows, err := a.db.QueryContext(ctx, account.SelectAll)
	if err != nil {
		return nil, err
	}

	var accs []bank.Account
	for rows.Next() {
		var acc bank.Account
		if sErr := rows.Scan(&acc.ID, &acc.Balance, &acc.Email); sErr != nil {
			err = errors.Join(err, sErr)
		}
		accs = append(accs, acc)
	}
	if err != nil {
		return nil, err
	}

	return accs, nil
}

func (a *accountRepo) Create(ctx context.Context, acc bank.Account) (bank.Account, error) {
	var id int64
	if err := a.db.QueryRowContext(ctx, account.CreateAccount, acc.Email, acc.Balance).Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return bank.Account{}, errs.ErrConflict
			}
		}
		return bank.Account{}, err
	}
	return bank.Account{
		ID:      id,
		Email:   acc.Email,
		Balance: acc.Balance,
	}, nil
}

func (a *accountRepo) GetByID(ctx context.Context, args bank.GetAccountArgs) (bank.Account, error) {
	res := a.db.QueryRowContext(ctx, account.GetByID, args.ID)
	var account bank.Account
	if err := res.Scan(&account.ID, &account.Balance, &account.Email); err != nil {
		return account, err
	}

	return account, nil
}

func (a *accountRepo) Update(ctx context.Context, req bank.UpdateAccountReq) (bank.Account, error) {
	panic("not implemented") // TODO: Implement
}

func (a *accountRepo) Delete(ctx context.Context, req bank.DeleteAccountReq) error {
	panic("not implemented") // TODO: Implement
}

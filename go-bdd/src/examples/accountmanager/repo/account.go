package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/accountmanager"
	"github.com/Tigh-Gherr/presentations/go-bdd/src/examples/accountmanager/repo/account"
)

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) accountmanager.AccountStore {
	return &accountRepo{
		db: db,
	}
}

func (a *accountRepo) List(ctx context.Context) ([]accountmanager.Account, error) {
	rows, err := a.db.QueryContext(ctx, account.SelectAll)
	if err != nil {
		return nil, err
	}

	var accs []accountmanager.Account
	for rows.Next() {
		var acc accountmanager.Account
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

func (a *accountRepo) Create(ctx context.Context, acc accountmanager.Account) (accountmanager.Account, error) {
	var id int64
	if err := a.db.QueryRowContext(ctx, account.CreateAccount, acc.Email, acc.Balance).Scan(&id); err != nil {
		return accountmanager.Account{}, err
	}
	return accountmanager.Account{
		ID:      id,
		Email:   acc.Email,
		Balance: acc.Balance,
	}, nil
}

func (a *accountRepo) GetByID(ctx context.Context, args accountmanager.GetAccountArgs) (accountmanager.Account, error) {
	res := a.db.QueryRowContext(ctx, account.GetByID, args.ID)
	var account accountmanager.Account
	if err := res.Scan(&account.ID, &account.Balance, &account.Email); err != nil {
		return account, err
	}

	return account, nil
}

func (a *accountRepo) Update(ctx context.Context, req accountmanager.UpdateAccountReq) (accountmanager.Account, error) {
	panic("not implemented") // TODO: Implement
}

func (a *accountRepo) Delete(ctx context.Context, req accountmanager.DeleteAccountReq) error {
	panic("not implemented") // TODO: Implement
}

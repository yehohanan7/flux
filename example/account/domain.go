package account

import (
	"github.com/golang/glog"
	"github.com/pborman/uuid"
	"github.com/yehohanan7/cqrs"
)

type Account struct {
	cqrs.Aggregate
	Id      string
	Balance int
}

func (account *Account) HandleNewAccount(event AccountCreated) {
	if account.Id != "" {
		glog.Error("attempt to create duplicate account")
		return
	}
	account.Id = event.AccountId
	account.Balance = event.Balance
}

func (account *Account) HandleCredit(event AccountCredited) {
	account.Balance = account.Balance + event.Amount
}

func (account *Account) HandleDebit(event AccountDebited) {
	account.Balance = account.Balance - event.Amount
}

func (account *Account) Credit(amount int) string {
	tId := uuid.New()
	account.Update(AccountCredited{tId, amount})
	return tId
}

func (account *Account) Debit(amount int) string {
	tId := uuid.New()
	account.Update(AccountDebited{tId, amount})
	return tId
}

func GetAccount(id string, store cqrs.EventStore) Account {
	account := new(Account)
	account.Aggregate = cqrs.NewAggregate(id, account, store)
	return *account
}

func MakeAccount(balance int, store cqrs.EventStore) Account {
	id := uuid.New()
	account := new(Account)
	account.Aggregate = cqrs.NewAggregate(id, account, store)
	account.Update(AccountCreated{id, balance})
	account.Save()
	return *account
}

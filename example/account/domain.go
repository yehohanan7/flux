package account

import (
	"fmt"

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
	fmt.Println("handling account created...")
	if account.Id != "" {
		glog.Error("attempt to create duplicate account")
		return
	}
	account.Id = event.AccountId
	account.Balance = event.Balance
}

func MakeAccount(balance int, store cqrs.EventStore) Account {
	id := uuid.New()
	account := new(Account)
	account.Aggregate = cqrs.NewAggregate(id, account, store)
	account.Update(AccountCreated{id, balance})
	account.Save()
	return *account
}

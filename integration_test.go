package cqrs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Account struct {
	Aggregate
	Id      string
	Balance int
}

type AccountCredited struct {
	Amount int
}

type AccountDebited struct {
	Amount int
}

func (acc *Account) HandleAccountCredited(event AccountCredited) {
	acc.Balance = acc.Balance + event.Amount
}

func (acc *Account) HandleAccountDebited(event AccountDebited) {
	acc.Balance = acc.Balance - event.Amount
}

func (acc *Account) Credit(amount int) {
	acc.Update(AccountCredited{amount})
}

func (acc *Account) Debit(amount int) {
	acc.Update(AccountDebited{amount})
}

func TestSaveAggregate(t *testing.T) {
	repo := NewInMemoryRepository()

	accountId := "account-id"
	account := new(Account)
	account.Aggregate = NewAggregate(accountId, account, repo)

	fmt.Println(account.Aggregate.handlers)

	account.Credit(5)
	account.Credit(10)
	account.Debit(1)

	assert.Equal(t, 14, account.Balance)
}
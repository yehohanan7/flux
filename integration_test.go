package cqrs

import (
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
	store := NewInMemoryEventStore()
	accountId := "account-id"
	account := new(Account)
	account.Aggregate = NewAggregate(accountId, account, store)
	account.Credit(5)
	account.Credit(10)
	account.Debit(1)

	account.Save()

	assert.Equal(t, 14, account.Balance)
}

func TestUpdateExistingAggregate(t *testing.T) {
	store := NewInMemoryEventStore()
	accountId := "account-id"
	existingAccount := new(Account)
	existingAccount.Aggregate = NewAggregate(accountId, existingAccount, store)
	existingAccount.Credit(5)

	existingAccount.Save()

	account := new(Account)
	account.Aggregate = NewAggregate(accountId, account, store)
	assert.Equal(t, 5, account.Balance)
}

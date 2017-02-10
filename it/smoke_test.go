package it

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yehohanan7/cqrs"
)

type Account struct {
	cqrs.Aggregate
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

type TransactionCountProjection struct {
	cqrs.Projection
	Count int
}

func (projection *TransactionCountProjection) HandleCredit(event AccountCredited) {
	projection.Count++
}

func (projection *TransactionCountProjection) HandleDebit(event AccountDebited) {
	projection.Count++
}

func TestAggregate(t *testing.T) {
	store := cqrs.NewEventStore()

	accountId := "account-id"
	existingAccount := new(Account)
	existingAccount.Aggregate = cqrs.NewAggregate(accountId, existingAccount, store)

	existingAccount.Credit(5)
	existingAccount.Credit(10)
	existingAccount.Debit(1)
	existingAccount.Save()

	account := new(Account)
	account.Aggregate = cqrs.NewAggregate(accountId, account, store)
	assert.Equal(t, 14, account.Balance)

	transactions := new(TransactionCountProjection)
	transactions.Projection = cqrs.NewProjection(accountId, transactions, store)
	assert.Equal(t, 3, transactions.Count)

}

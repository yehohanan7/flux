package account

import (
	"github.com/golang/glog"
	"github.com/yehohanan7/flux"
	"github.com/yehohanan7/flux/cqrs"
)

type Summary struct {
	Id      string
	Balance int
}

type AccountEventConsumer struct {
	cqrs.EventConsumer
	accounts map[string]int
}

func (ac *AccountEventConsumer) HandleNewAccount(event AccountCreated) {
	glog.Info("Handling new account")
}

func (ac *AccountEventConsumer) HandleCredits(event AccountCredited) {
	glog.Info("Handling credits ", event)
	ac.accounts[event.AccountId] += event.Amount
}

func (ac *AccountEventConsumer) HandleDebits(event AccountDebited) {
	glog.Info("Handling debits", event)
	ac.accounts[event.AccountId] -= event.Amount
}

func (ac *AccountEventConsumer) GetSummary(id string) Summary {
	glog.Info(ac.accounts)
	return Summary{id, ac.accounts[id]}
}

func NewAccontSummaryConsumer(url string) *AccountEventConsumer {
	consumer := new(AccountEventConsumer)
	consumer.accounts = make(map[string]int)
	consumer.EventConsumer = flux.NewEventConsumer(url, consumer)
	return consumer
}

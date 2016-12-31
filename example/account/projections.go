package account

import (
	"github.com/golang/glog"
	"github.com/yehohanan7/cqrs"
)

type AccountSummary struct {
	cqrs.Projection
	Id             string `json:"id"`
	CurrentBalance int    `json:"current_balance"`
}

func (p *AccountSummary) HandleNewAccount(event AccountCreated) {
	glog.Info("Handling new account: ", event)
	p.Id = event.AccountId
	p.CurrentBalance = event.Balance
}

func (p *AccountSummary) HandleCredits(event AccountCredited) {
	glog.Info("Handling credit: ", event)
	p.CurrentBalance = p.CurrentBalance + event.Amount
}

func (p *AccountSummary) HandleDebits(event AccountDebited) {
	glog.Info("Handling debit: ", event)
	p.CurrentBalance = p.CurrentBalance - event.Amount
}

func ProjectAccountSummary(accountId string) AccountSummary {
	summary := new(AccountSummary)
	summary.Projection = cqrs.NewProjection(accountId, summary, store)
	return *summary
}

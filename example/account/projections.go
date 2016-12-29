package account

import (
	"github.com/yehohanan7/cqrs"
)

type AccountSummary struct {
	cqrs.Projection
	Id             string `json:"id"`
	CurrentBalance int    `json:"current_balance"`
}

func (p *AccountSummary) HandleNewAccount(event AccountCreated) {
	p.Id = event.AccountId
	p.CurrentBalance = event.Balance
}

func ProjectAccountSummary(accountId string) AccountSummary {
	summary := new(AccountSummary)
	summary.Projection = cqrs.NewProjection(accountId, summary, store)
	return *summary
}

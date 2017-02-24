package account

type AccountSummary struct {
	Id             string `json:"id"`
	CurrentBalance int    `json:"current_balance"`
}

func (p *AccountSummary) HandleNewAccount(event AccountCreated) {

}

func (p *AccountSummary) HandleCredits(event AccountCredited) {

}

func (p *AccountSummary) HandleDebits(event AccountDebited) {

}

func ProjectAccountSummary(accountId string) AccountSummary {
	summary := new(AccountSummary)
	return *summary
}

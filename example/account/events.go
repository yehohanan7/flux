package account

type AccountCreated struct {
	AccountId string
	Balance   int
}

type AccountCredited struct {
	AccountId string
	Amount    int
}

type AccountDebited struct {
	AccountId string
	Amount    int
}

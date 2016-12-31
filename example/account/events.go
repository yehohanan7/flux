package account

type AccountCreated struct {
	AccountId string
	Balance   int
}

type AccountCredited struct {
	TransactionId string
	Amount        int
}

type AccountDebited struct {
	TransactionId string
	Amount        int
}

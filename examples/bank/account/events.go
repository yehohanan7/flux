package account

type AccountCreated struct {
	AccountId string `json:"account_id"`
	Balance   int    `json:"balance"`
}

type AccountCredited struct {
	AccountId     string `json:"account_id"`
	TransactionId string `json:"transaction_id"`
	Amount        int    `json:"amount"`
}

type AccountDebited struct {
	AccountId     string `json:"account_id"`
	TransactionId string `json:"transaction_id"`
	Amount        int    `json:"amount"`
}

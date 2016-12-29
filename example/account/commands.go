package account

type Command interface {
	Execute() (interface{}, error)
}

type CreateAccountCommand struct {
	OpeningBalance int `json:"opening_balance"`
}

type CreditAccountCommand struct {
	Amount int `json:"amount"`
}

type DebitAccountCommand struct {
	Amount int `json:"amount"`
}

func (command *CreateAccountCommand) Execute() (interface{}, error) {
	return MakeAccount(command.OpeningBalance, store).Id, nil
}

func (command *CreditAccountCommand) Execute() (interface{}, error) {
	return nil, nil
}

func (command *DebitAccountCommand) Execute() (interface{}, error) {
	return nil, nil
}

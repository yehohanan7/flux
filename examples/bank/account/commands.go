package account

import (
	"github.com/golang/glog"
)

type Command interface {
	Execute() (interface{}, error)
}

type CreateAccountCommand struct {
	OpeningBalance int `json:"opening_balance"`
}

type CreditAccountCommand struct {
	AccountId string `json:"account_id"`
	Amount    int    `json:"amount"`
}

type DebitAccountCommand struct {
	AccountId string `json:"account_id"`
	Amount    int    `json:"amount"`
}

func (command *CreateAccountCommand) Execute() (interface{}, error) {
	glog.Info("Executing create account command:", command)
	id := MakeAccount(command.OpeningBalance, store).Id
	glog.Infof("account with id %s created\n", id)
	return id, nil
}

func (command *CreditAccountCommand) Execute() (interface{}, error) {
	glog.Info("Executing credit account command:", command)
	account := GetAccount(command.AccountId, store)
	transId := account.Credit(command.Amount)
	account.Save()
	glog.Infof("Account %s credited with %v dollars\n", command.AccountId, command.Amount)
	return transId, nil
}

func (command *DebitAccountCommand) Execute() (interface{}, error) {
	glog.Info("Executing debit account command:", command)
	account := GetAccount(command.AccountId, store)
	transId := account.Debit(command.Amount)
	account.Save()
	glog.Infof("Account %s debited with %v dollars\n", command.AccountId, command.Amount)
	return transId, nil
}

package account

import (
	"github.com/golang/glog"
	"github.com/yehohanan7/flux"
)

type Summary struct {
	Id      string
	Balance int
}

type AccountRepository struct {
	accounts map[string]int
}

func (r *AccountRepository) Get(id string) Summary {
	return Summary{id, r.accounts[id]}
}

func NewAccountSummaryRepository(url string) *AccountRepository {
	repo := new(AccountRepository)
	events := []interface{}{AccountCreated{}, AccountCreated{}, AccountDebited{}}
	store := flux.NewOffsetStore()
	consumer := flux.NewEventConsumer(url, events, store)
	go func() {
		eventCh := make(chan interface{})
		go consumer.Start(eventCh, nil)
		for {
			event := <-eventCh
			switch event.(type) {
			case AccountCreated:
				glog.Info("Handling new account")
				e := event.(AccountCreated)
				repo.accounts[e.AccountId] = e.Balance
			case AccountCredited:
				glog.Info("Handling credits ", event)
				e := event.(AccountCreated)
				repo.accounts[e.AccountId] += e.Balance
			case AccountDebited:
				glog.Info("Handling debits", event)
				e := event.(AccountCreated)
				repo.accounts[e.AccountId] -= e.Balance
			}
		}
	}()
	return repo
}

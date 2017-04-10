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
	repo := &AccountRepository{make(map[string]int)}
	events := []interface{}{AccountCreated{}, AccountCredited{}, AccountDebited{}}
	store := flux.NewMemoryOffsetStore()
	consumer := flux.NewEventConsumer(url, events, store)
	go func() {
		eventCh := make(chan interface{})
		glog.Info("Starting consumer...")
		go consumer.Start(eventCh)
		glog.Info("Consumer started")
		for {
			event := <-eventCh
			glog.Info("recieved event ", event)
			switch event.(type) {
			case AccountCreated:
				glog.Info("Handling new account")
				e := event.(AccountCreated)
				repo.accounts[e.AccountId] = e.Balance
			case AccountCredited:
				glog.Info("Handling credits ", event)
				e := event.(AccountCredited)
				repo.accounts[e.AccountId] += e.Amount
			case AccountDebited:
				glog.Info("Handling debits", event)
				e := event.(AccountDebited)
				repo.accounts[e.AccountId] -= e.Amount
			}
		}
	}()
	return repo
}

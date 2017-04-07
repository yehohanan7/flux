# Flux [![Build Status](https://travis-ci.org/yehohanan7/flux.svg)](https://travis-ci.org/yehohanan7/flux?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/yehohanan7/flux)](https://goreportcard.com/report/github.com/yehohanan7/flux)
![logo](http://www.logogala.com/images/uploads/gallery/octopus.png)


# Introduction
"There is nothing called state. There are events and the story we tell about what it means."

If you want to try CQRS & DDD, you would have a service which accepts commands on an aggregate and publish messages to a messaging system like kafka, and you will have various services consuming these messages and building views/read model from these messages.

But if you feel you just need a simple applicaiton without the hassle of horizontal scalability of kafka, then Flux is the answer for you!


## Aggregate
Flux suggests that you use one service per [Aggregate](http://serviceorientation.com/soaglossary/entity_service), which accepts commands and publishes events.

This is how you would define an aggregate in Flux:

```go
import(
  "github.com/yehohanan7/flux"
	"github.com/yehohanan7/flux/cqrs"
)

//Account is an aggregate
type Account struct {
	cqrs.Aggregate
	Id      string
	Balance int
}

//Initialize the aggregate
acc := new(Account)
acc.Aggregate = cqrs.NewAggregate("account-id", acc, flux.NewMemoryStore())
```

The last argument is an implementation of EventStore interface, there are 2 implementations at the moment an inmemory one and a boltdb implementation
```go
store := flux.NewBoltStore("path/to/database")
```

Once you have the aggregate initialized, you can execute commands on it which will emit events
```go
//My event
type AccountCredited struct {
	Amount int
}

//Command
func (acc *Account) Credit(amount int) {
	acc.Update(AccountCredited{amount})
}

//Event handler to allow you to update the state of the aggregate
func (acc *Account) HandleAccountCredited(event AccountCredited) {
	acc.Balance = acc.Balance + event.Amount
}

```
*Note that you should have the handler method prefixed with the name "Handle"*


ok, now you need to ensure remote services gets the events from this aggregate, which is pretty simple if you have a mux router in your application.

```go
	router.HandleFunc("/events", flux.FeedHandler(store))
	router.HandleFunc("/events/{id}", flux.FeedHandler(store))
```

## Read models

okay, now you need to build a read model from the events published by the aggregate.

```go
//stores the offset to know where to consumer from after a restart
offsetStore := flux.NewOffsetStore()
consumer := flux.NewEventConsumer(url, events, offsetStore)
eventCh := make(chan interface{})

//Start consuming
go consumer.Start(eventCh)

//Fetch the events and build your read models
for {
	event := <-eventCh
	switch event.(type) {
  case AccountCredited:
    fmt.Println(event.(AccountCredited))
  }
}
```

# Flux [![Build Status](https://travis-ci.org/yehohanan7/flux.svg)](https://travis-ci.org/yehohanan7/flux?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/yehohanan7/flux)](https://goreportcard.com/report/github.com/yehohanan7/flux)
![logo](http://www.logogala.com/images/uploads/gallery/octopus.png)


# Introduction
*"There is nothing called state. There are events and the story we tell about what it means."*

Flux allows you to quickly build an application in CQRS way without the hassle of a messaging system like RabbitMQ or Kafka inbetween your command and read model services.

It's a good practice to have one command service per Aggregate (as per DDD terminology) and various read model/view services. the command service stores the events that are emited by each command and expose the same as a json feed for the consumers (read model services) to consume in regular intervals allowing you to easily decouple commands and read model services.


## Aggregate
Flux suggests that you use one service per [Aggregate](https://martinfowler.com/bliki/DDD_Aggregate.html), which accepts commands and publishes events.

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

The last argument is an EventStore, which provides an implementation to store and retrieve events - there are 2 implementations at the moment an inmemory one and a boltdb implementation
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

//Event handler to allow you to update the state of the aggregate as a result of a command
func (acc *Account) HandleAccountCredited(event AccountCredited) {
	acc.Balance = acc.Balance + event.Amount
}


//Execute command
acc.Credit(100)
acc.Credit(150)
acc.Save()
```
**Note that you should have the handler method prefixed with the name "Handle"**


ok, now you need to ensure remote services gets the events from this aggregate, which is pretty simple if you have a mux router in your application.

```go
router.HandleFunc("/events", flux.FeedHandler(store))
router.HandleFunc("/events/{id}", flux.FeedHandler(store))
```

Attaching the feed handler to your router publishes the events of the aggregate as a feed, which would look like below

```json
{
  "description": "event feed",
  "events": [
    {
      "event_id": "47d074c3-ba83-40e2-8720-804b73a202b9",
      "url": "http://localhost:3000/events/47d074c3-ba83-40e2-8720-804b73a202b9",
      "aggregate_name": "*account.Account",
      "aggregate_version": 0,
      "event_type": "account.AccountCreated",
      "created": "Fri Apr  7 15:19:18 2017"
    },
    {
      "event_id": "174a40b6-104a-4be5-a352-4e61b524d3dc",
      "url": "http://localhost:3000/events/174a40b6-104a-4be5-a352-4e61b524d3dc",
      "aggregate_name": "*account.Account",
      "aggregate_version": 1,
      "event_type": "account.AccountCredited",
      "created": "Fri Apr  7 15:19:27 2017"
    }
  ]
}
```

## Read models
you don't have to process the json feed to build your read model, you can just start a Flux event consumer and provide it with a list of events you are interested it and you will get back all the events over a channel

```go
//stores the offset to know where to consumer from after a restart
offsetStore := flux.NewOffsetStore()
consumer := flux.NewEventConsumer("http://entityservicehost:port/events", []interface{}{AccountCredited{}}, offsetStore)
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

You could pause,resume & stop the consumer
```go
consumer.Pause()
consumer.Resume()
consumer.Stop()
```

## Read model storage
As you notice, flux doesn't support storage of your read models. once you get the events, you could store the events/state in any storage system. however, flux will provide a default storage for storing read models in future releases


## Sample application
There is a simple example application [here](https://github.com/yehohanan7/flux/tree/master/examples/bank) if you would like to refer

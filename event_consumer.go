package cqrs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/golang/glog"
)

const defaultOffset = 0

//Consumes events from the command component
type EventConsumer struct {
	url          string
	handlerClass interface{}
	handlers     handlermap
}

//Send event to the consumer
func (consumer *EventConsumer) send(event Event) {
	payload := event.Payload
	if handler, ok := consumer.handlers[reflect.TypeOf(payload)]; ok {
		handler(consumer.handlerClass, payload)
	}
}

func fetchJson(url string, data interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		glog.Error("Error while getting ", url, err)
		return err
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Error("Error while reading the data ", err)
		return err
	}

	err = json.Unmarshal(body, data)

	if err != nil {
		glog.Error("Error while decoding data ", err)
		return err
	}

	return nil
}

func getFeed(url string) (JsonEventFeed, error) {
	var feed = new(JsonEventFeed)
	err := fetchJson(url, feed)
	if err != nil {
		return nil, err
	}
	return *feed, nil
}

func getEvent(url string) (EventEntry, error) {
	var event *EventEntry
	err := fetchJson(url, event)
	if err != nil {
		return nil, err
	}
	return *event, nil
}

func (consumer *EventConsumer) Start() error {
	feed, err := getFeed(consumer.url)
	if err != nil {
		return err
	}

	for _, entry := range feed.Events {
		event := getEvent(entry.Url)
		if handler, ok := consumer.handlers[reflect.TypeOf(payload)]; ok {

		}
	}
	return nil
}

func (consumer *EventConsumer) Stop() error {
	return nil
}

//Create new consumer
func NewEventConsumer(url string, handlerClass interface{}, store OffsetStore) EventConsumer {
	return EventConsumer{url, handlerClass, buildHandlerMap(handlerClass)}
}

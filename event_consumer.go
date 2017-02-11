package cqrs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/golang/glog"
)

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

func fetchJsonInto(url string, data interface{}) error {
	var body []byte
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

func (consumer *EventConsumer) getEventFeed() (JsonEventFeed, error) {
	var feed = new(JsonEventFeed)
	err := fetchJsonInto(consumer.url, feed)
	if err != nil {
		return *feed, err
	}
	return *feed, nil
}

func (consumer *EventConsumer) getEvent(entry EventEntry) (interface{}, error) {

	for eventType, _ := range consumer.handlers {
		if eventType.String() == entry.EventType {
			event := reflect.New(eventType)
			err := fetchJsonInto(entry.Url, event.Interface())
			if err == nil {
				return event.Elem().Interface(), err
			}
		}
	}

	return nil, nil
}

func (consumer *EventConsumer) Start() error {
	feed, err := consumer.getEventFeed()
	if err != nil {
		return err
	}

	for _, entry := range feed.Events {
		event, err := consumer.getEvent(entry)
		if err != nil {
			panic(err)
		}
		if handler, ok := consumer.handlers[reflect.TypeOf(event)]; ok {
			handler(consumer.handlerClass, event)
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

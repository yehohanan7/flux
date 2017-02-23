package consumer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
	. "github.com/yehohanan7/flux/feed"
)

//Consumes events from the command component
type JsonEventConsumer struct {
	url          string
	handlerClass interface{}
	handlers     Handlers
}

//Send event to the consumer
func (consumer *JsonEventConsumer) send(event Event) {
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

func (consumer *JsonEventConsumer) getEventFeed() (JsonEventFeed, error) {
	var feed = new(JsonEventFeed)
	err := fetchJsonInto(consumer.url, feed)
	if err != nil {
		return *feed, err
	}
	return *feed, nil
}

func (consumer *JsonEventConsumer) getEvent(entry EventEntry) (interface{}, error) {

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

func (consumer *JsonEventConsumer) Start() error {
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

func (consumer *JsonEventConsumer) Stop() error {
	return nil
}

//Create new consumer
func NewEventConsumer(url string, handlerClass interface{}, store OffsetStore) *JsonEventConsumer {
	return &JsonEventConsumer{url, handlerClass, NewHandlers(handlerClass)}
}

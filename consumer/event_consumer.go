package consumer

import (
	"reflect"

	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
	. "github.com/yehohanan7/flux/feed"
	"github.com/yehohanan7/flux/utils"
)

type SimpleConsumer struct {
	url    string
	events map[string]reflect.Type
	store  OffsetStore
}

func (c *SimpleConsumer) consume(entry EventEntry) interface{} {
	if e, ok := c.events[entry.EventType]; ok {
		event := reflect.New(e).Interface()
		err := utils.HttpGetJson(entry.Url, event)
		if err != nil {
			glog.Error("error while getting event ", entry.Url, err)
			return nil
		}
		return reflect.ValueOf(event).Elem().Interface()
	}
	return nil
}

func (c *SimpleConsumer) Start(eventCh, stopCh chan interface{}) {
	var feed = new(JsonEventFeed)
	err := utils.HttpGetJson(c.url, feed)
	if err != nil {
		close(eventCh)
		return
	}

	for _, entry := range feed.Events {
		if event := c.consume(entry); event != nil {
			eventCh <- event
		}
	}
}

func New(url string, events []interface{}, store OffsetStore) *SimpleConsumer {
	m := make(map[string]reflect.Type)
	for _, e := range events {
		m[reflect.TypeOf(e).String()] = reflect.TypeOf(e)
	}
	return &SimpleConsumer{url, m, store}
}

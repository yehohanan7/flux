package consumer

import (
	"reflect"

	"github.com/golang/glog"
	. "github.com/yehohanan7/flux/cqrs"
	. "github.com/yehohanan7/flux/feed"
	"github.com/yehohanan7/flux/utils"
)

type SimpleConsumer struct {
	url   string
	em    map[string]reflect.Type
	store OffsetStore
}

func fetch(em map[string]reflect.Type, entry EventEntry) interface{} {
	if e, ok := em[entry.EventType]; ok {
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
		if event := fetch(c.em, entry); event != nil {
			eventCh <- event
		}
	}
}

func eventMap(events []interface{}) map[string]reflect.Type {
	m := make(map[string]reflect.Type)
	for _, e := range events {
		m[reflect.TypeOf(e).String()] = reflect.TypeOf(e)
	}
	return m
}

func New(url string, events []interface{}, store OffsetStore) *SimpleConsumer {
	return &SimpleConsumer{url, eventMap(events), store}
}

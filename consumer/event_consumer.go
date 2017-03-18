package consumer

import (
	"reflect"
	"strconv"
	"time"

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

func (c *SimpleConsumer) Start(eventCh, stopCh chan interface{}) error {
	for _ = range time.Tick(5 * time.Second) {
		var feed = new(JsonEventFeed)
		offset, err := c.store.GetLastOffset()
		if err != nil {
			glog.Error("Error while getting last offset ", err)
			return err
		}

		glog.Info("Fetching events from offset ", offset)
		err = utils.HttpGetJson(c.url+"?offset="+strconv.Itoa(offset), feed)
		if err != nil {
			close(eventCh)
			return err
		}

		for _, entry := range feed.Events {
			if event := fetch(c.em, entry); event != nil {
				eventCh <- event
			}
		}
		c.store.SaveOffset(offset + len(feed.Events))
	}
	return nil
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

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

func eventMap(events []interface{}) map[string]reflect.Type {
	m := make(map[string]reflect.Type)
	for _, e := range events {
		m[reflect.TypeOf(e).String()] = reflect.TypeOf(e)
	}
	return m
}

type SimpleConsumer struct {
	url             string
	em              map[string]reflect.Type
	store           OffsetStore
	pollInterval    time.Duration
	paused, stopped bool
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

func getFeed(url string, offset int) (JsonEventFeed, error) {
	glog.Info("Fetching events from offset ", offset)
	var feed = new(JsonEventFeed)
	err := utils.HttpGetJson(url+"?offset="+strconv.Itoa(offset), feed)
	return *feed, err
}

func (c *SimpleConsumer) Start(eventCh chan interface{}) error {
	for _ = range time.Tick(c.pollInterval) {
		if c.stopped {
			close(eventCh)
			return nil
		}
		if !c.paused {
			offset, _ := c.store.GetLastOffset()
			feed, err := getFeed(c.url, offset)
			if err != nil {
				glog.Error("Error while getting feed ", err)
				return err
			}
			for _, entry := range feed.Events {
				if event := fetch(c.em, entry); event != nil {
					eventCh <- event
				}
			}
			c.store.SaveOffset(offset + len(feed.Events))
		}
	}
	return nil
}

func (c *SimpleConsumer) Pause() {
	glog.Info("Pausing consumer...")
	c.paused = true
}

func (c *SimpleConsumer) Resume() {
	glog.Info("Resuming consumer...")
	c.paused = false
}

func (c *SimpleConsumer) Stop() {
	glog.Info("Stopping consumer...")
	c.stopped = true
}

func New(url string, events []interface{}, store OffsetStore, pollInterval time.Duration) *SimpleConsumer {
	return &SimpleConsumer{url, eventMap(events), store, pollInterval, false, false}
}

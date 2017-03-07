package consumer

import (
	. "github.com/yehohanan7/flux/cqrs"
	. "github.com/yehohanan7/flux/feed"
	"github.com/yehohanan7/flux/utils"
)

type SimpleConsumer struct {
	url    string
	events []interface{}
	store  OffsetStore
}

func (c *SimpleConsumer) Start(eventCh, stopCh chan interface{}) {
	for {
		var feed = new(JsonEventFeed)
		err := utils.HttpGetJson(c.url, feed)
		if err != nil {
			close(eventCh)
			return
		}
		eventCh <- *feed
	}
}

func New(url string, events []interface{}, store OffsetStore) *SimpleConsumer {
	return &SimpleConsumer{url, events, store}
}

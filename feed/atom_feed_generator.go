package feed

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	. "github.com/yehohanan7/cqrs/cqrs"
)

type AtomFeedGenerator struct {
}

func (_ AtomFeedGenerator) ContentType() string {
	return "text/xml"
}

func (_ AtomFeedGenerator) Generate(url, description string, events []Event) string {
	feed := &feeds.Feed{
		Title:       description,
		Link:        &feeds.Link{Href: url},
		Description: description,
	}

	feed.Items = make([]*feeds.Item, 0)

	for _, event := range events {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          event.Id,
			Title:       event.AggregateName,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/%s", url, event.Id)},
			Author:      &feeds.Author{Name: "cqrs", Email: "cqrs@cqrs.org"},
			Description: reflect.TypeOf(event.Payload).String(),
			Created:     time.Now(),
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		glog.Error("error while building atom feed", err)
	}
	return atom
}

package cqrs

import (
	"encoding/json"
	"net/http"
	"time"

	"reflect"

	"fmt"

	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

const DEFAULT_PAGE_SIZE = 20

var pageSize = DEFAULT_PAGE_SIZE

func generateFeed(url string, store EventStore, offset int) string {

	feed := &feeds.Feed{
		Title:       "event feeds",
		Link:        &feeds.Link{Href: url},
		Description: "All events on this service",
	}

	feed.Items = make([]*feeds.Item, 0)

	for _, event := range store.GetAllEventsFrom(offset, pageSize) {
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

func getFeed(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(generateFeed(r.URL.RequestURI(), store, offset)))
	}
}

func getEventById(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)
		id := vars["id"]
		json.NewEncoder(w).Encode(store.GetEvent(id))
	}
}

func EventFeed(router *mux.Router, store EventStore, eventsPerPage ...int) {
	if len(eventsPerPage) > 1 {
		panic("invalid number of arguments")
	}

	if len(eventsPerPage) == 1 && eventsPerPage[0] <= 0 {
		panic("invalid events per page")
	}

	if len(eventsPerPage) == 1 {
		pageSize = eventsPerPage[0]
	}

	router.HandleFunc("/events", getFeed(store)).Methods("GET")
	router.HandleFunc("/events/{id}", getEventById(store)).Methods("GET")
}

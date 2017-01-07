package cqrs

import (
	"encoding/json"
	"net/http"
	"time"

	"reflect"

	"fmt"

	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

func generateFeed(url string, store EventStore) string {

	feed := &feeds.Feed{
		Title:       "event feeds",
		Link:        &feeds.Link{Href: url},
		Description: "All events on this service",
	}

	feed.Items = make([]*feeds.Item, 0)

	for _, event := range store.GetAllEvents() {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          event.Id,
			Title:       event.AggregateName,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/%s", url, event.Id)},
			Author:      &feeds.Author{Name: "cqrs", Email: "cqrs@localhost.com"},
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
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(generateFeed(r.URL.RequestURI(), store)))
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

func EventFeed(router *mux.Router, store EventStore) {
	router.HandleFunc("/events", getFeed(store)).Methods("GET")
	router.HandleFunc("/events/{id}", getEventById(store)).Methods("GET")
}

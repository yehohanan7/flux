package cqrs

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

const DEFAULT_PAGE_SIZE = 20

var pageSize = DEFAULT_PAGE_SIZE

type FeedGenerator interface {
	Generate(string, string, []Event) string
	ContentType() string
}

var generators = map[string]FeedGenerator{
	"json": JsonFeedGenerator{},
	"atom": AtomFeedGenerator{},
}

func getFeed(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		generator := generators["json"]
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		events := store.GetAllEventsFrom(offset, pageSize)
		feedType := r.FormValue("format")

		if _, ok := generators[feedType]; ok {
			generator = generators[feedType]
		}

		w.Header().Set("Content-Type", generator.ContentType())
		w.Write([]byte(generator.Generate(absoluteUrl(r), "event feed", events)))
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

//Exposes events as atom feed
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

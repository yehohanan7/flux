package feed

import (
	"encoding/json"
	"net/http"
	"strconv"

	"strings"

	. "github.com/yehohanan7/flux/cqrs"
	. "github.com/yehohanan7/flux/utils"
)

const DEFAULT_PAGE_SIZE = 20

var generator = JsonFeedGenerator{}

func events(w http.ResponseWriter, r *http.Request, store EventStore) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	metas := store.GetEventMetaDataFrom(offset, DEFAULT_PAGE_SIZE)
	w.Header().Set("Content-Type", "application/json")
	w.Write(generator.Generate(GetAbsoluteUrl(r), "event feed", metas))
}

func event(w http.ResponseWriter, r *http.Request, store EventStore, id string) {
	event := store.GetEvent(id)
	json.NewEncoder(w).Encode(event.Payload)
}

func getEventId(path string) string {
	xs := strings.Split(path, "/")
	if len(xs) == 3 && len(xs[2]) > 0 {
		return xs[2]
	}
	return ""
}

//Exposes events as atom feed
func GetFeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		events(w, r, store)
	}
}

//Get event by event id
func GetEventHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if id := getEventId(r.URL.Path); len(id) > 0 {
			event(w, r, store, id)
		}
	}
}

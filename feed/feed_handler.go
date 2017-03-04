package feed

import (
	"encoding/json"
	"fmt"
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
	events := store.GetAllEventsFrom(offset, DEFAULT_PAGE_SIZE)
	w.Header().Set("Content-Type", "application/json")
	w.Write(generator.Generate(GetAbsoluteUrl(r), "event feed", events))
}

func event(w http.ResponseWriter, r *http.Request, store EventStore, id string) {
	event := store.GetEvent(id)
	fmt.Println("fetched event: ", event)
	json.NewEncoder(w).Encode(event)
}

//Exposes events as atom feed
func FeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.URL.Path == "/events" {
			events(w, r, store)
			return
		}

		ps := strings.Split(r.URL.Path, "/")
		if len(ps) == 3 && len(ps[2]) > 0 {
			event(w, r, store, ps[2])
		}
	}
}

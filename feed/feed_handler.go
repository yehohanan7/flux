package feed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

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
	json.NewEncoder(w).Encode(store.GetEvent(id))
}

//Exposes events as atom feed
func FeedHandler(store EventStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile("/events/*")
		defer r.Body.Close()
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/events" {
			events(w, r, store)
		}
		if ids := re.FindStringSubmatch(r.URL.Path); len(ids) > 1 {
			fmt.Println(ids)
			event(w, r, store, ids[0])
		}
	}
}

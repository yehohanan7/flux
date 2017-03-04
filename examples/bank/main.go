package main

import (
	"encoding/json"
	"log"
	"net/http"

	"flag"

	"github.com/gorilla/mux"
	"github.com/yehohanan7/flux"
	. "github.com/yehohanan7/flux/examples/bank/account"
)

var (
	store    = flux.NewEventStore()
	consumer = NewAccontSummaryConsumer("http://localhost:3000/events")
)

func init() {
	InitAccounts(store)
	err := consumer.Start()
	if err != nil {
		panic(err)
	}
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	json.NewEncoder(w).Encode(consumer.GetSummary(id))
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	command := new(CreateAccountCommand)
	json.NewDecoder(r.Body).Decode(command)
	response, _ := command.Execute()
	json.NewEncoder(w).Encode(response)
}

func CreditAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	command := new(CreditAccountCommand)
	json.NewDecoder(r.Body).Decode(command)
	command.AccountId = mux.Vars(r)["id"]
	response, _ := command.Execute()
	json.NewEncoder(w).Encode(response)
}

func DebitAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	command := new(DebitAccountCommand)
	json.NewDecoder(r.Body).Decode(command)
	command.AccountId = mux.Vars(r)["id"]
	response, _ := command.Execute()
	json.NewEncoder(w).Encode(response)
}

func main() {
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/accounts", CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}/summary", GetSummary).Methods("GET")
	router.HandleFunc("/accounts/{id}/credit", CreditAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}/debit", DebitAccount).Methods("POST")
	router.HandleFunc("/events/*", flux.FeedHandler(store))
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", router))
}

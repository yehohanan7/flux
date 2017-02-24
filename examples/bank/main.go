package main

import (
	"encoding/json"
	"log"
	"net/http"

	"flag"

	"github.com/gorilla/mux"
	. "github.com/yehohanan7/bank/account"
	"github.com/yehohanan7/cqrs"
)

var store cqrs.EventStore

func init() {
	store = cqrs.NewEventStore()
	InitAccounts(store)
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	json.NewEncoder(w).Encode(ProjectAccountSummary(id))
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
	cqrs.EventFeed(router, store)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", router))
}

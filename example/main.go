package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yehohanan7/cqrs"
	. "github.com/yehohanan7/cqrs/example/account"
)

func init() {
	InitAccounts(cqrs.NewInMemoryEventStore())
}

func ExecuteCommand(w http.ResponseWriter, r *http.Request, command Command) {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(command)
	if err != nil {
		fmt.Println(err)
		return
	}
	response, _ := command.Execute()
	json.NewEncoder(w).Encode(response)
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	json.NewEncoder(w).Encode(ProjectAccountSummary(id))
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	ExecuteCommand(w, r, new(CreateAccountCommand))
}

func CreditAccount(w http.ResponseWriter, r *http.Request) {
	ExecuteCommand(w, r, new(CreditAccountCommand))
}

func DebitAccount(w http.ResponseWriter, r *http.Request) {
	ExecuteCommand(w, r, new(DebitAccountCommand))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}/summary", GetSummary).Methods("GET")
	router.HandleFunc("/accounts/{id}/credit", CreditAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}/debit", DebitAccount).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", router))
}

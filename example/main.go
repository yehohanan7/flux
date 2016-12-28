package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateAccount() {

}

func Router() (r *mux.Router) {
	r = mux.NewRouter()
	r.HandleFunc("/{id}/balance", GetBalance).Methods("GET")
	r.HandleFunc("/{id}/credit", CreditAccount).Methods("POST")
	r.HandleFunc("/{id}/debit", DebitAccount).Methods("POST")
	return
}

func main() {
	router := Router()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", router))
}

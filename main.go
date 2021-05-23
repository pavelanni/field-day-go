package main

import (
	"field-day/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

type Visitor struct {
	Name      string
	Callsign  string
	FirstTime bool
	Nfarl     bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "fieldday_dev"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/new", usersC.New).Methods("GET")
	r.HandleFunc("/new", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}

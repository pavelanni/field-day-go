package main

import (
	"field-day/controllers"
	"field-day/models"
	"fmt"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	us.AutoMigrate()

	staticDir := "/static"

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/new", usersC.New).Methods("GET")
	r.HandleFunc("/new", usersC.Create).Methods("POST")
	r.PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	http.ListenAndServe(":3000", r)
}

package main

import (
	"field-day/controllers"
	"field-day/models"
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
	dbfile = "field-day.db"
)

func main() {
	us, err := models.NewUserService(dbfile)
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
	r.Handle("/confirmation", staticC.Confirmation).Methods("GET")
	r.HandleFunc("/new", usersC.New).Methods("GET")
	r.HandleFunc("/new", usersC.Create).Methods("POST")
	r.PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	http.ListenAndServe(":3000", r)
}

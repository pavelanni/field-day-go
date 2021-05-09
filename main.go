package main

import (
	"field-day/controllers"
	"field-day/views"
	"net/http"

	"github.com/gorilla/mux"
)

type Visitor struct {
	Name      string
	Callsign  string
	FirstTime bool
	Nfarl     bool
}

var (
	homeView     *views.View
	contactView  *views.View
	listView     *views.View
	notfoundView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(notfoundView.Render(w, nil))
}

func list(w http.ResponseWriter, r *http.Request) {
	data := []Visitor{{"Pavel Anni", "AC4PA", false, true},
		{"Jim Stafford", "W4QO", false, true},
		{"Gene Shablygin", "W3UA", true, false}}

	w.Header().Set("Content-Type", "text/html")
	must(listView.Render(w, data))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	staticC := controllers.NewStatic()
	notfoundView = views.NewView("bootstrap", "static/notfound")
	listView = views.NewView("bootstrap", "list")
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.HandleFunc("/list", list).Methods("GET")
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/new", usersC.New).Methods("GET")
	r.HandleFunc("/new", usersC.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.ListenAndServe(":3000", r)
}

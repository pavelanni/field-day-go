package main

import (
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
	signupView   *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil))
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
	homeView = views.NewView("bootstrap",
		"views/home.go.html")
	contactView = views.NewView("bootstrap",
		"views/contact.go.html")
	notfoundView = views.NewView("bootstrap",
		"views/notfound.go.html")
	listView = views.NewView("bootstrap",
		"views/list.go.html")
	signupView = views.NewView("bootstrap",
		"views/new.go.html")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/list", list)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/new", signup)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.ListenAndServe(":3000", r)
}

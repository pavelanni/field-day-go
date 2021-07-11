package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

type Visitor struct {
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
	Callsign  string `schema:"callsign"`
	Email     string `schema:"email"`
	Nfarl     bool   `schema:"nfarl"`
	Contactme bool   `schema:"contactme"`
	Youth     bool   `schema:"youth"`
	Firsttime bool   `schema:"firsttime"`
}

var templateDir = "templates"
var staticDir = "static"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		templateDir + "/bootstrap-refresh.go.html",
		templateDir + "/header.go.html",
		templateDir + "/home.go.html",
		templateDir + "/footer.go.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func confirmHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		templateDir + "/bootstrap-refresh.go.html",
		templateDir + "/header.go.html",
		templateDir + "/confirmation.go.html",
		templateDir + "/footer.go.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		files := []string{
			templateDir + "/bootstrap.go.html",
			templateDir + "/header.go.html",
			templateDir + "/new.go.html",
			templateDir + "/footer.go.html",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)

		if err != nil {
			log.Fatal(err)
		}
		return
	}
	// If POST
	v := Visitor{}
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(&v, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, v)
}

func main() {
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/confirmation", confirmHandler)
	http.ListenAndServe(":3000", nil)
}

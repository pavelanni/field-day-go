package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

// Visitor contains information about Field Day visitors
type Visitor struct {
	gorm.Model
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
var dbFile string

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

func listHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		templateDir + "/bootstrap.go.html",
		templateDir + "/header.go.html",
		templateDir + "/list.go.html",
		templateDir + "/footer.go.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	visitors, err := listVisitors(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, visitors)
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
		log.Fatal(err)
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(&v, r.PostForm); err != nil {
		log.Fatal(err)
	}
	if err := saveVisitor(v); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}

// Run opens the database and starts the server
// Using Run func allows us to test it later (we can't test the main func)
// I learned it from Elliot of TutorialEdge.net
func Run() error {
	err := createTable(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/confirmation", confirmHandler)
	http.HandleFunc("/list", listHandler)
	log.Println("Listening on port 3000")
	return http.ListenAndServe(":3000", nil)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Please provide a database file")
	}
	dbFile = os.Args[1]
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/schema"
	"github.com/pavelanni/field-day-go/visitorstore"
)

var templateDir = "templates"
var staticDir = "static"

type Server struct {
	store *visitorstore.VisitorStore
}

func NewServer(dbFile string) (*Server, error) {
	store, err := visitorstore.NewVisitorStore(dbFile)
	if err != nil {
		return nil, err
	}
	return &Server{store}, nil
}

func (s *Server) Run() error {
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/new", s.NewVisitorHandler)
	http.HandleFunc("/confirmation", confirmHandler)
	http.HandleFunc("/list", s.ListHandler)
	log.Println("Listening on port 3000")
	return http.ListenAndServe(":3000", nil)
}

// Handlers
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

func (s *Server) ListHandler(w http.ResponseWriter, r *http.Request) {
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

	visitors, err := s.store.ListVisitors()
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, visitors)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) NewVisitorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		files := []string{
			templateDir + "/bootstrap.go.html",
			templateDir + "/header.go.html",
			templateDir + "/new.go.html",
			templateDir + "/footer.go.html",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
		return
	}
	// If POST
	v := visitorstore.Visitor{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parsing error; please return to the previous page", http.StatusInternalServerError)
		return
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(&v, r.PostForm); err != nil {
		http.Error(w, "Form decoding error; please return to the previous page", http.StatusInternalServerError)
		return
	}
	if err := s.store.SaveVisitor(v); err != nil {
		http.Error(w, "Visitor saving error; please return to the previous page", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Please provide a database file")
	}
	dbFile := os.Args[1]

	server, err := NewServer(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

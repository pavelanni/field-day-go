package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/schema"
	"github.com/pavelanni/field-day-go/visitorstore"
)

var templateDir = "templates"
var staticDir = "static"
var port = "3000"
var thisYear = "2025"

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/css/* static/js/* static/NFARL_FD_2025.png static/nfarlLogoTransparentBackground_medium.gif
var staticFS embed.FS

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
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/new", s.NewVisitorHandler)
	http.HandleFunc("/confirmation", confirmHandler)
	http.HandleFunc("/list", s.ListHandler)
	log.Println("Listening on port " + port)
	return http.ListenAndServe(":"+port, nil)
}

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		templateDir + "/bootstrap-refresh.go.html",
		templateDir + "/header.go.html",
		templateDir + "/home.go.html",
		templateDir + "/footer.go.html",
	}
	tmpl, err := template.ParseFS(templatesFS, files...)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]any{"Year": thisYear}
	err = tmpl.Execute(w, data)

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
	tmpl, err := template.ParseFS(templatesFS, files...)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]any{"Year": thisYear}
	err = tmpl.Execute(w, data)

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
	tmpl, err := template.ParseFS(templatesFS, files...)
	if err != nil {
		log.Fatal(err)
	}

	visitors, err := s.store.ListVisitors()
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]any{"Visitors": visitors, "Year": thisYear}
	err = tmpl.Execute(w, data)
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
		tmpl, err := template.ParseFS(templatesFS, files...)
		if err != nil {
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			return
		}
		totalVisitors, err := s.store.TotalVisitors()
		if err != nil {
			log.Fatal(err)
		}
		data := map[string]any{"Year": thisYear, "CurrentVisitor": totalVisitors + 1}
		err = tmpl.Execute(w, data)

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

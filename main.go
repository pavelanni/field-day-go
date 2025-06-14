package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/schema"
	"github.com/pavelanni/field-day-go/morse"
	"github.com/pavelanni/field-day-go/visitorstore"
)

var templateDir = "templates"
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

func (s *Server) run() error {
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/new", s.newVisitorHandler)
	http.HandleFunc("/confirmation", confirmHandler)
	http.HandleFunc("/list", s.listHandler)
	http.HandleFunc("/morse-audio", morseAudioHandler)
	log.Println("Listening on port " + port)
	return http.ListenAndServe(":"+port, nil)
}

// handlers
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
	callsign := r.URL.Query().Get("callsign")
	name := r.URL.Query().Get("name")
	data := map[string]any{"Year": thisYear, "Callsign": callsign, "Name": name}
	err = tmpl.Execute(w, data)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) listHandler(w http.ResponseWriter, r *http.Request) {
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

// Morse code is now played in the browser via generated WAV files served by the backend.
// The morse package's Play method is retained for possible future CLI/desktop use.
func (s *Server) newVisitorHandler(w http.ResponseWriter, r *http.Request) {
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
	msg := v.Callsign
	if msg != "" {
		msg = msg + "  73"
	} else {
		msg = "73"
	}
	http.Redirect(w, r, "/confirmation?callsign="+url.QueryEscape(msg)+"&name="+url.QueryEscape(v.FirstName), http.StatusSeeOther)
}

// Handler to serve Morse code audio as WAV
func morseAudioHandler(w http.ResponseWriter, r *http.Request) {
	callsign := r.URL.Query().Get("callsign")
	if callsign == "" {
		http.Error(w, "Missing callsign parameter", http.StatusBadRequest)
		return
	}
	audioData, err := morse.GenerateWav(callsign)
	if err != nil {
		http.Error(w, "Failed to generate audio", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Disposition", "inline; filename=\"morse.wav\"")
	w.Write(audioData)
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

	err = server.run()
	if err != nil {
		log.Fatal(err)
	}
}

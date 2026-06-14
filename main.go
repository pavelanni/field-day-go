package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/pavelanni/field-day-go/memberlookup"
	"github.com/pavelanni/field-day-go/morse"
	"github.com/pavelanni/field-day-go/visitorstore"
	"github.com/spf13/pflag"
)

var templateDir = "templates"
var defaultPort = "3000"
var thisYear = "2026"

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/css/* static/js/* static/NFARL_FD_2026.png static/nfarlLogoTransparentBackground_medium.gif
var staticFS embed.FS

type Server struct {
	store     *visitorstore.VisitorStore
	members   *memberlookup.Lookup
	templates map[string]*template.Template
	year      string
}

func NewServer(dbFile, membersFile, year string) (*Server, error) {
	store, err := visitorstore.NewVisitorStore(dbFile)
	if err != nil {
		return nil, err
	}

	var members *memberlookup.Lookup
	if membersFile != "" {
		members, err = memberlookup.LoadCSV(membersFile)
		if err != nil {
			return nil, fmt.Errorf("loading members CSV: %w", err)
		}
		log.Printf("Loaded %d club members from %s", members.Len(), membersFile)
	} else {
		log.Println("No members CSV specified; member auto-fill disabled")
	}

	s := &Server{store: store, members: members, year: year}
	if err := s.initTemplates(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) initTemplates() error {
	s.templates = map[string]*template.Template{}

	// Pre-parse template bundles used by handlers
	// home
	if tmpl, err := template.ParseFS(templatesFS,
		templateDir+"/tailwind-refresh.go.html",
		templateDir+"/header.go.html",
		templateDir+"/home.go.html",
		templateDir+"/footer.go.html",
	); err != nil {
		return err
	} else {
		s.templates["home"] = tmpl
	}

	// confirmation
	if tmpl, err := template.ParseFS(templatesFS,
		templateDir+"/tailwind-refresh.go.html",
		templateDir+"/header.go.html",
		templateDir+"/confirmation.go.html",
		templateDir+"/footer.go.html",
	); err != nil {
		return err
	} else {
		s.templates["confirm"] = tmpl
	}

	// list
	if tmpl, err := template.ParseFS(templatesFS,
		templateDir+"/tailwind.go.html",
		templateDir+"/header.go.html",
		templateDir+"/list.go.html",
		templateDir+"/footer.go.html",
	); err != nil {
		return err
	} else {
		s.templates["list"] = tmpl
	}

	// new visitor
	if tmpl, err := template.ParseFS(templatesFS,
		templateDir+"/tailwind.go.html",
		templateDir+"/header.go.html",
		templateDir+"/new.go.html",
		templateDir+"/footer.go.html",
	); err != nil {
		return err
	} else {
		s.templates["new"] = tmpl
	}

	// privacy
	if tmpl, err := template.ParseFS(templatesFS,
		templateDir+"/tailwind-refresh-timeout.go.html",
		templateDir+"/header.go.html",
		templateDir+"/privacy.go.html",
		templateDir+"/footer.go.html",
	); err != nil {
		return err
	} else {
		s.templates["privacy"] = tmpl
	}

	return nil
}

func (s *Server) run(addr, port string) error {
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	mux.HandleFunc("/", s.homeHandler)
	mux.HandleFunc("/new", s.newVisitorHandler)
	mux.HandleFunc("/confirmation", s.confirmHandler)
	mux.HandleFunc("/list", s.listHandler)
	mux.HandleFunc("/morse-audio", s.morseAudioHandler)
	mux.HandleFunc("/privacy", s.privacyHandler)
	mux.HandleFunc("/healthz", s.healthHandler)
	mux.HandleFunc("/api/visitor-count", s.visitorCountHandler)
	mux.HandleFunc("/member-lookup", s.memberLookupHandler)

	srv := &http.Server{
		Addr:              addr + ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	log.Println("Listening on " + srv.Addr)
	return srv.ListenAndServe()
}

// handlers
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := s.templates["home"]
	data := map[string]any{"Year": s.year}
	if err := tmpl.ExecuteTemplate(w, "tailwind-refresh", data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Printf("homeHandler template error: %v", err)
		return
	}
}

func (s *Server) confirmHandler(w http.ResponseWriter, r *http.Request) {
	callsign := r.URL.Query().Get("callsign")
	name := r.URL.Query().Get("name")
	data := map[string]any{"Year": s.year, "Callsign": callsign, "Name": name}
	tmpl := s.templates["confirm"]
	if err := tmpl.ExecuteTemplate(w, "tailwind-refresh", data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Printf("confirmHandler template error: %v", err)
		return
	}
}

func (s *Server) listHandler(w http.ResponseWriter, r *http.Request) {
	visitors, err := s.store.ListVisitors()
	if err != nil {
		http.Error(w, "Failed to load visitors", http.StatusInternalServerError)
		log.Printf("listHandler store error: %v", err)
		return
	}
	data := map[string]any{"Visitors": visitors, "Year": s.year}
	tmpl := s.templates["list"]
	if err := tmpl.ExecuteTemplate(w, "tailwind", data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Printf("listHandler template error: %v", err)
		return
	}
}

// Morse code is now played in the browser via generated WAV files served by the backend.
// The morse package's Play method is retained for possible future CLI/desktop use.
func (s *Server) newVisitorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := s.templates["new"]
		totalVisitors, err := s.store.TotalVisitors()
		if err != nil {
			http.Error(w, "Failed to get totals", http.StatusInternalServerError)
			log.Printf("newVisitorHandler totals error: %v", err)
			return
		}
		data := map[string]any{"Year": s.year, "CurrentVisitor": totalVisitors + 1, "TotalVisitors": totalVisitors}
		if err := tmpl.ExecuteTemplate(w, "tailwind", data); err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			log.Printf("newVisitorHandler template error: %v", err)
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

	// Normalize fields (Validate normalizes callsign and email)
	v.FirstName = strings.TrimSpace(v.FirstName)
	v.LastName = strings.TrimSpace(v.LastName)

	// Validate fields
	if err := v.Validate(); err != nil {
		if valErr, ok := err.(*visitorstore.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			tmpl := s.templates["new"]
			totalVisitors, err := s.store.TotalVisitors()
			if err != nil {
				log.Printf("newVisitorHandler totals error on validation re-render: %v", err)
			}
			data := map[string]any{
				"Year":           s.year,
				"CurrentVisitor": totalVisitors + 1,
				"TotalVisitors":  totalVisitors,
				"Error":          valErr.Error(),
			}
			if err := tmpl.ExecuteTemplate(w, "tailwind", data); err != nil {
				http.Error(w, "Template execution error", http.StatusInternalServerError)
				log.Printf("newVisitorHandler template error: %v", err)
			}
			return
		}
	}

	if err := s.store.SaveVisitor(v); err != nil {
		http.Error(w, "Visitor saving error; please return to the previous page", http.StatusInternalServerError)
		log.Printf("newVisitorHandler save error: %v", err)
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
func (s *Server) morseAudioHandler(w http.ResponseWriter, r *http.Request) {
	callsign := r.URL.Query().Get("callsign")
	if callsign == "" {
		http.Error(w, "Missing callsign parameter", http.StatusBadRequest)
		return
	}
	audioData, err := morse.GenerateWav(callsign)
	if err != nil {
		http.Error(w, "Failed to generate audio", http.StatusInternalServerError)
		log.Printf("morseAudioHandler error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Disposition", "inline; filename=\"morse.wav\"")
	w.Write(audioData)
}

func (s *Server) privacyHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{"Year": s.year}
	tmpl := s.templates["privacy"]
	if err := tmpl.ExecuteTemplate(w, "tailwind-refresh-timeout", data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Printf("privacyHandler template error: %v", err)
		return
	}
}

// healthHandler returns 200 and checks DB connectivity
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.store.HealthCheck(); err != nil {
		http.Error(w, "db not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) visitorCountHandler(w http.ResponseWriter, r *http.Request) {
	count, err := s.store.TotalVisitors()
	if err != nil {
		http.Error(w, "failed to get visitor count", http.StatusInternalServerError)
		log.Printf("visitorCountHandler error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"year": s.year, "count": count}); err != nil {
		log.Printf("visitorCountHandler encode error: %v", err)
	}
}

func (s *Server) memberLookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if s.members == nil {
		_, _ = w.Write([]byte(`{}`))
		return
	}

	callsign := r.URL.Query().Get("callsign")
	if callsign == "" {
		_, _ = w.Write([]byte(`{}`))
		return
	}

	member, ok := s.members.Lookup(callsign)
	if !ok {
		_, _ = w.Write([]byte(`{}`))
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"first_name": member.FirstName,
		"last_name":  member.LastName,
		"email":      member.Email,
	}); err != nil {
		log.Printf("memberLookupHandler encode error: %v", err)
	}
}

func main() {
	var (
		flagDB      string
		flagPort    string
		flagAddr    string
		flagYear    string
		flagMembers string
	)

	// Defaults can be overridden by environment variables
	flagDB = os.Getenv("FD_DB")
	if flagDB == "" && len(os.Args) > 1 {
		// Backward compatibility: first positional arg is DB file
		flagDB = os.Args[1]
	}
	flagPort = os.Getenv("FD_PORT")
	if flagPort == "" {
		flagPort = defaultPort
	}
	flagAddr = os.Getenv("FD_ADDR")
	if flagAddr == "" {
		flagAddr = "0.0.0.0"
	}
	flagYear = os.Getenv("FD_YEAR")
	if flagYear == "" {
		flagYear = thisYear
	}
	flagMembers = os.Getenv("FD_MEMBERS")

	pflag.StringVar(&flagDB, "db", flagDB, "Path to SQLite DB file (required)")
	pflag.StringVar(&flagPort, "port", flagPort, "Port to listen on")
	pflag.StringVar(&flagAddr, "addr", flagAddr, "Bind address")
	pflag.StringVar(&flagYear, "year", flagYear, "Event year for templates")
	pflag.StringVar(&flagMembers, "members", flagMembers, "Path to club members CSV file")
	pflag.Parse()

	if flagDB == "" {
		log.Fatal("database file not provided: use --db or FD_DB or positional arg")
	}

	server, err := NewServer(flagDB, flagMembers, flagYear)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.run(flagAddr, flagPort); err != nil {
		log.Fatal(err)
	}
}

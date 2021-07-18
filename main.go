package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
var dbFile = "fd2021.db"

func saveVisitor(v Visitor) error {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return err
	}
	result := db.Create(&v)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

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

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var visitors []Visitor
	result := db.Find(&visitors)
	if result.Error != nil {
		log.Fatal(result.Error)
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

func createTable(dbFile string) error {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) { // create the table
		db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Visitor{})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
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
	http.ListenAndServe(":3000", nil)
}

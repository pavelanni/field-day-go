package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

type Visitor struct {
	FirstName string `form:"firstname"`
	LastName  string `form:"lastname"`
	Callsign  string `form:"callsign"`
	Email     string `form:"email"`
	Nfarl     bool   `form:"nfarl"`
	Contactme bool   `form:"contactme"`
	Youth     bool   `form:"youth"`
	Firsttime bool   `form:"firsttime"`
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
	r.ParseForm()
	// This is my own attempt to implement something like gorilla.schema
	// It was good to use it to learn more about reflect package
	// but probably I should use the existing gorilla.schema
	v := Visitor{}
	vt := reflect.TypeOf(v)
	vp := reflect.ValueOf(&v)
	vpe := vp.Elem()
	for i := 0; i < vt.NumField(); i++ {
		if vpe.Field(i).Kind() == reflect.String {
			vpe.Field(i).SetString(r.PostForm[vt.Field(i).Tag.Get("form")][0])
		}
		if vpe.Field(i).Kind() == reflect.Bool {
			if len(r.PostForm[vt.Field(i).Tag.Get("form")]) == 1 {
				vpe.Field(i).SetBool(true)
			} else {
				vpe.Field(i).SetBool(false)
			}
		}
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

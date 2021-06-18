package controllers

import (
	"field-day/models"
	"field-day/views"
	"net/http"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

type SignupForm struct {
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
	Callsign  string `schema:"callsign"`
	Email     string `schema:"email"`
	Nfarl     bool   `schema:"nfarl"`
	Contactme bool   `schema:"contactme"`
	Youth     bool   `schema:"youth"`
	Firsttime bool   `schema:"firsttime"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	type Alert struct {
		Level   string
		Message string
	}
	type Data struct {
		Alert *Alert
		Yield interface{}
	}

	alert := Alert{
		Level:   "success",
		Message: "Successfully registered a new visitor!",
	}
	data := Data{
		Alert: &alert,
		Yield: "NFARL Field Day 2021",
	}
	if err := u.NewView.Render(w, data); err != nil {
		panic(err)
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Callsign:  form.Callsign,
		Email:     form.Email,
		Nfarl:     form.Nfarl,
		Contactme: form.Contactme,
		Firsttime: form.Firsttime,
		Youth:     form.Youth,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}

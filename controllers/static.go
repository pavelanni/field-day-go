package controllers

import "field-day/views"

func NewStatic() *Static {
	return &Static{
		Home:     views.NewView("bootstrap", "static/home"),
		Contact:  views.NewView("bootstrap", "static/contact"),
		Notfound: views.NewView("bootstrap", "static/notfound"),
	}
}

type Static struct {
	Home     *views.View
	Contact  *views.View
	Notfound *views.View
}

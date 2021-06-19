package controllers

import "field-day/views"

func NewStatic() *Static {
	return &Static{
		Home:         views.NewView("bootstrap-refresh", "static/home"),
		Contact:      views.NewView("bootstrap", "static/contact"),
		Confirmation: views.NewView("bootstrap-refresh", "static/confirmation"),
	}
}

type Static struct {
	Home         *views.View
	Contact      *views.View
	Confirmation *views.View
}

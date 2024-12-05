package controllers

import (
	"html/template"
	"net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		http.Error(w, "Cannot load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
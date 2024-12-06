package pages

import (
	"html/template"
	"net/http"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Cannot load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
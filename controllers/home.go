package controllers

import (
	"html/template"
	"http-server/config"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session-name") // 使用 session 儲存登入狀態
	username, ok := session.Values["username"].(string)

	data := map[string]interface{}{
		"IsLoggedIn": ok,
		"Username":   username,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Cannot load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
package controllers

import (
	"html/template"
	"http-server/controllers/applications"
	"net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(applications.UsernameContextKey).(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // 未找到用戶名時重定向到登入頁面
		return
	}

	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		http.Error(w, "Cannot load template", http.StatusInternalServerError)
		return
	}

	// username, ok := r.Context().Value(applications.UsernameContextKey).(string)
	// if !ok {
	//     http.Redirect(w, r, "/login", http.StatusSeeOther) // 未找到用戶名時重定向到登入頁面
	//     return
	// }

	// 將 username 傳遞到模板
	data := map[string]string{
		"Username": username,
	}

	tmpl.Execute(w, data)
}
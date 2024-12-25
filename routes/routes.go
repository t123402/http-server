package routes

import (
	"http-server/controllers" // 匯入控制器
	"http-server/controllers/applications"
	"http-server/controllers/pages"
	"net/http"
)

func Routes() {
	http.HandleFunc("/", controllers.HomeHandler)
	http.Handle("/about", applications.Authenticate(http.HandlerFunc(controllers.AboutHandler)))
	itemRoutes()
	authRoutes()
}

func itemRoutes() {
	http.HandleFunc("/items", pages.ItemsPageHandler)
	http.HandleFunc("/api/items", applications.GetItemsHandler)
	http.HandleFunc("/api/items/add", applications.AddItemHandler)
	http.HandleFunc("/api/items/delete/", applications.DeleteItemHandler)
	http.HandleFunc("/api/items/update/", applications.UpdateItemHandler)
}

func authRoutes() {
	http.HandleFunc("/login", pages.LoginPageHandler)
	http.HandleFunc("/auth/register", applications.RegisterHandler)
	http.HandleFunc("/auth/login", applications.LoginHandler)
	http.HandleFunc("/auth/logout", applications.LogoutHandler)
	http.HandleFunc("/auth/profile/", applications.ProfileHandler)
	http.HandleFunc("/auth/profile/update/", applications.UpdateProfileHandler)
	http.Handle("/auth/me", applications.Authenticate(http.HandlerFunc(applications.MeHandler)))
}

package routes

import (
	"http-server/controllers" // 匯入控制器
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/", controllers.HomeHandler)
	http.HandleFunc("/about", controllers.AboutHandler)
	http.HandleFunc("/items", controllers.ItemsHandler)
	http.HandleFunc("/api/items", controllers.GetItemsHandler)
	http.HandleFunc("/api/items/add", controllers.AddItemHandler)
	http.HandleFunc("/api/items/delete/", controllers.DeleteItemHandler)
	http.HandleFunc("/api/items/update/", controllers.UpdateItemHandler)
}
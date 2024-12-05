package main

import (
	"fmt"
	"http-server/database"
	"http-server/routes" // 匯入路由設定
	"net/http"
)

func main() {
	// 初始化資料庫
	database.InitDB()
	
	// 註冊路由
	routes.RegisterRoutes()

	// 啟動伺服器，監聽在 8080 埠號
	fmt.Println("伺服器啟動，監聽在 http://localhost")
	http.ListenAndServe(":8080", nil)
}
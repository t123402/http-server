package main

import (
	"fmt"
	"http-server/config"
	"http-server/routes" // 匯入路由設定
	"net/http"
)

func main() {
	// 初始化全局配置
	config.InitConfig()

	// 註冊路由
	routes.Routes()

	// 啟動伺服器，監聽在 8080 埠號
	fmt.Println("伺服器啟動，監聽在 http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
package routes

import (
	// 匯入控制器
	"http-server/controllers"
	"net/http"
	"os"
	"path/filepath"
)

func Routes() {
	// 獲取當前工作目錄
	workingDir, _ := os.Getwd()
	// 拼接出靜態文件的絕對路徑，假設靜態文件存放在 "public" 資料夾內
	staticPath := filepath.Join(workingDir, "public")

	// 註冊一個處理所有 HTTP 請求的處理函數
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 拼接請求的路徑到靜態文件目錄，得到請求的完整文件路徑
		filePath := filepath.Join(staticPath, r.URL.Path)

		// 檢查請求的文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// 如果文件不存在，返回 React 的 "index.html"
			// 這是為了支持 React 的前端路由，確保所有未知路徑都由前端處理
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
		} else {
			// 如果文件存在，直接返回靜態文件
			http.ServeFile(w, r, filePath)
		}
	})

	itemRoutes()
	authRoutes()
}

func itemRoutes() {
	http.HandleFunc("/api/items", controllers.GetItemsHandler)
	http.HandleFunc("/api/items/add", controllers.AddItemHandler)
	http.HandleFunc("/api/items/delete/", controllers.DeleteItemHandler)
	http.HandleFunc("/api/items/update/", controllers.UpdateItemHandler)
}

func authRoutes() {
	http.HandleFunc("/auth/register", controllers.RegisterHandler)
	http.HandleFunc("/auth/login", controllers.LoginHandler)
	http.HandleFunc("/auth/logout", controllers.LogoutHandler)
	http.HandleFunc("/auth/profile/", controllers.ProfileHandler)
	http.HandleFunc("/auth/profile/update/", controllers.UpdateProfileHandler)
	http.HandleFunc("/auth/change-password/", controllers.ChangePasswordHandler)
	http.Handle("/auth/me", controllers.Authenticate(http.HandlerFunc(controllers.MeHandler)))
}

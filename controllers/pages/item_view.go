package pages

import (
	"html/template"
	"net/http"
)

// ItemsHandler 處理渲染 items 頁面的請求
func ItemsPageHandler(w http.ResponseWriter, r *http.Request) {
	// 加載模板檔案
	tmpl, err := template.ParseFiles("templates/item.html")
	if err != nil {
		// 如果模板檔案無法加載，返回 HTTP 500 錯誤
		http.Error(w, "Cannot load template", http.StatusInternalServerError)
		return
	}
	// 渲染模板並輸出到回應
	tmpl.Execute(w, nil)
}
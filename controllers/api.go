package controllers

import (
	"encoding/json"
	"html/template"
	"http-server/database"
	"net/http"
	"strings"
)

// Item 是用來表示 items 資料表中的一個資料結構
type Item struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// ItemsHandler 處理渲染 items 頁面的請求
func ItemsHandler(w http.ResponseWriter, r *http.Request) {
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

// 新增資料
func AddItemHandler(w http.ResponseWriter, r *http.Request) {
	// 確認請求方法是否為 POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 解析請求體中的 JSON，並將其映射到 Item 結構
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 將資料插入到資料庫
	query := "INSERT INTO items (value) VALUES (?)"
	_, err = database.DB.Exec(query, item.Value)
	if err != nil {
		// 插入資料失敗，返回 HTTP 500 錯誤
		http.Error(w, "Failed to insert item", http.StatusInternalServerError)
		return
	}

	// 返回 HTTP 201 Created，表示新增成功
	w.WriteHeader(http.StatusCreated)
}

// 查詢資料
func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	// 從資料庫中查詢所有資料
	rows, err := database.DB.Query("SELECT id, value FROM items")
	if err != nil {
		// 查詢失敗，返回 HTTP 500 錯誤
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // 確保查詢結果的資源在使用完後被正確釋放

	var items []Item
	// 遍歷查詢結果，將每一行映射到 Item 結構，並添加到 items 列表中
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Value)
		if err != nil {
			// 映射失敗，返回 HTTP 500 錯誤
			http.Error(w, "Failed to parse item", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	// 將 items 列表以 JSON 格式返回給用戶
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// 刪除資料
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	// 確認請求方法是否為 DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 從 URL 中提取 ID，例如 /api/items/delete/1 中提取到 "1"
	id := strings.TrimPrefix(r.URL.Path, "/api/items/delete/")
	if id == "" {
		// 如果未提供 ID，返回 HTTP 400 錯誤
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	// 執行刪除操作
	query := "DELETE FROM items WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		// 刪除失敗，返回 HTTP 500 錯誤
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	// 返回 HTTP 204 No Content，表示刪除成功且無內容返回
	w.WriteHeader(http.StatusNoContent)
}

// 修改資料
func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	// 確認請求方法是否為 PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 從 URL 中提取 ID，例如 /api/items/update/1 中提取到 "1"
	id := strings.TrimPrefix(r.URL.Path, "/api/items/update/")
	if id == "" {
		// 如果未提供 ID，返回 HTTP 400 錯誤
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	// 解析請求體中的 JSON，並將其映射到 Item 結構
	type Item struct {
		Value string `json:"value"`
	}
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		// 請求體格式錯誤，返回 HTTP 400 錯誤
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 更新資料庫中的資料
	query := "UPDATE items SET value = ? WHERE id = ?"
	_, err := database.DB.Exec(query, item.Value, id)
	if err != nil {
		// 更新失敗，返回 HTTP 500 錯誤
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	// 返回 HTTP 200 OK，表示更新成功
	w.WriteHeader(http.StatusOK)
}
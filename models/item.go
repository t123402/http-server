package models

import (
	"http-server/database"
)

// Item 是用來表示 items 資料表中的一個資料結構
type Item struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

// GetAllItems 查詢所有 items
func GetAllItems() ([]Item, error) {
	// 從資料庫中查詢所有資料
	rows, err := database.DB.Query("SELECT id, value FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 確保查詢結果的資源在使用完後被正確釋放

	var items []Item
	// 遍歷查詢結果，將每一行映射到 Item 結構，並添加到 items 列表中
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Value); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// AddItem 新增一個 item
func AddItem(value string) error {
	query := "INSERT INTO items (value) VALUES (?)"
	_, err := database.DB.Exec(query, value)
	return err
}

// DeleteItem 根據 ID 刪除 item
func DeleteItem(id string) error {
	query := "DELETE FROM items WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	return err
}

// UpdateItem 根據 ID 更新 item
func UpdateItem(id, value string) error {
	query := "UPDATE items SET value = ? WHERE id = ?"
	_, err := database.DB.Exec(query, value, id)
	return err
}
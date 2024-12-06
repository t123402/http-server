package models

import (
	"http-server/database"
)

// User 表示 users 資料表中的一條記錄
type User struct {
	ID           int
	Username     string
	PasswordHash string
	RoleID       string
}

// AddUser 新增用戶
func AddUser(username, passwordHash, roleID string) error {
	query := "INSERT INTO users (username, password_hash, role_id) VALUES (?, ?, ?)"
	_, err := database.DB.Exec(query, username, passwordHash, roleID)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByUsername 根據用戶名查詢用戶
func GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, password_hash, role_id FROM users WHERE username = ?"
	row := database.DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RoleID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
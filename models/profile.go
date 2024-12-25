package models

import (
	"fmt"
	"http-server/database"
	"time"
)

type Profile struct {
	UserID    int
	Username  string
	Nickname  string
	Firstname string
	Lastname  string
	Email     string
	Gender    string
	Birthday  *time.Time
}

// AddProfile 新增用戶資訊
func AddProfile(username, nickname, firstname, lastname, email, gender string, birthday *time.Time) error {
	query := "INSERT INTO profiles (username, nickname, firstname, lastname, email, gender, birthday) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := database.DB.Exec(query, username, nickname, firstname, lastname, email, gender, birthday)
	if err != nil {
		return err
	}
	return nil
}

// GetProfileByUsername 根據用戶名查詢用戶資訊
func GetProfileByUsername(username string) (*Profile, error) {
	query := "SELECT user_id, username, nickname, firstname, lastname, email, gender, birthday FROM profiles WHERE username = ?"
	row := database.DB.QueryRow(query, username)

	var profile Profile
	err := row.Scan(&profile.UserID, &profile.Username, &profile.Nickname, &profile.Firstname, &profile.Lastname, &profile.Email, &profile.Gender, &profile.Birthday)
	if err != nil {
		fmt.Printf("Scan error: %v\n", err) // 輸出具體的錯誤
		return nil, err
	}
	return &profile, nil
}

// UpdateProfileByUsername 根據用戶名更新用戶資訊
func UpdateProfileByUsername(username, nickname, firstname, lastname, email, gender string, birthday *time.Time) error {
	query := `
		UPDATE profiles 
		SET nickname = ?, 
			firstname = ?, 
			lastname = ?, 
			email = ?, 
			gender = ?, 
			birthday = ? 
		WHERE username = ?
	`
	_, err := database.DB.Exec(query, nickname, firstname, lastname, email, gender, birthday, username)
	if err != nil {
		fmt.Printf("Update error: %v\n", err) // 輸出具體的錯誤
		return err
	}
	return nil
}

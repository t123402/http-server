package models

import (
	"http-server/database"
)

type Role struct {
	ID          int
	Name        string
	Description string
}

// GetRoleById 查詢用戶角色
func GetRoleById(id string) (*Role, error) {
	query := "SELECT id, name, description FROM roles WHERE id = ?"
	row := database.DB.QueryRow(query, id)

	var role Role
	err := row.Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

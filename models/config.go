package models

import (
	"http-server/database"
)

// ConfigEntry 表示 config 表中的一條記錄
type ConfigEntry struct {
	Key   string
	Value string
}

// GetAllConfigs 查詢所有配置
func GetAllConfigs() ([]ConfigEntry, error) {
	query := "SELECT `key`, value FROM config"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []ConfigEntry
	for rows.Next() {
		var config ConfigEntry
		if err := rows.Scan(&config.Key, &config.Value); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
package config

import (
	"fmt"
	"http-server/database"
	"http-server/models"
	"os"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// ConfigManager 是一個用於管理應用程序配置的結構體。
// configMap 用於存儲配置的 key-value 鍵值對，
// mu 是一個讀寫鎖，確保在多線程環境中對 configMap 的操作是安全的。
type ConfigManager struct {
	configMap map[string]string
	mu        sync.RWMutex // 讀寫鎖，用於保護 configMap
}

// instance 是 ConfigManager 的全局唯一實例，保證配置的單例性。
var instance *ConfigManager

// once 用於確保 LoadConfig() 中的初始化邏輯只執行一次。
var once sync.Once

var Store *sessions.CookieStore

// InitConfig 初始化所有配置
func InitConfig() {
	var err error

	// 加載 .env 文件
	err = godotenv.Load()
	if err != nil {
		// 如果 .env 文件無法加載，程序無法繼續執行
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	// 初始化資料庫
	database.InitDB()

	// 初始化全局 Session Store
	Store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,  // 設置存活時間（秒）
		HttpOnly: true,  // 禁止 JavaScript 訪問
	}

	err = LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to do LoadConfig: %v", err))
	}

	fmt.Println("Config initialized successfully")
}

// LoadConfig 從資料庫加載配置，並初始化全局單例。
// 如果已經加載過配置，則不會重複執行。
func LoadConfig() error {
	// once.Do 確保以下邏輯只執行一次，即使有多個線程同時調用 LoadConfig()。
	once.Do(func() {
		// 初始化 ConfigManager 實例
		instance = &ConfigManager{configMap: make(map[string]string)}

		// 從資料庫加載配置
		configs, err := models.GetAllConfigs()
		if err != nil {
			// 如果查詢資料庫失敗，清空實例並拋出錯誤
			instance = nil
			panic(fmt.Errorf("failed to load config: %v", err))
		}

		// 臨時存儲配置的鍵值對
		tempMap := make(map[string]string)
		for _, config := range configs {
			tempMap[config.Key] = config.Value
		}

		// 將臨時 map 賦值給 ConfigManager 實例的 configMap
		instance.configMap = tempMap
	})
	return nil
}

// GetProperty 根據給定的 key 獲取配置值。
// 返回值分為兩部分：1. 對應的 value；2. 是否存在該 key 的標誌（布爾值）。
func (c *ConfigManager) GetProperty(key string) (string, bool) {
	// 使用讀鎖保護 configMap，確保多線程下讀操作是安全的
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.configMap[key] // 查找配置 key
	return value, exists             // 返回結果和是否存在的標誌
}
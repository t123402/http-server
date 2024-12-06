package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL/MariaDB 驅動
)

// DB 是全局變量，用於存儲資料庫連線對象
// 通過 InitDB 初始化後，全局可以使用 DB 來進行資料庫操作
var DB *sql.DB

// InitDB 初始化資料庫連線
func InitDB() {
	var err error

	// 從環境變數中獲取資料庫連線字串（DSN）
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		panic("環境變數 DATABASE_DSN 未設定")
	}

	// 使用資料庫連線字串建立連線池
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("資料庫連線失敗: %v", err))
	}

	// 驗證資料庫連線是否成功
	// Ping 用於測試資料庫連線是否可用，這可以捕獲任何潛在的連線問題
	err = DB.Ping()
	if err != nil {
		panic(fmt.Sprintf("資料庫無法連線: %v", err))
	}

	fmt.Println("資料庫連線成功")
}
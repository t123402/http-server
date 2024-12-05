# HTTP Server

這是一個用 Go 語言開發的簡單 HTTP Server。它可以處理靜態頁面顯示和 API 請求，並與資料庫進行交互，實現基本的 CRUD 操作。

# 使用 Nginx 配置 HTTPS

通過 Nginx 設置 HTTPS，保證數據傳輸的安全性。後端服務監聽在 8080 埠，Nginx 作為反向代理，處理 80 埠的 HTTP 請求並重定向到 443 埠的 HTTPS 請求。免費的 SSL/TLS 證書由 Let's Encrypt 提供，確保應用能安全地運行於網路上。
ps 目前沒有網域 沒辦法設置 https 只能先用公有 IP 來訪問 http server

# 使用 .env 文件設置資料庫連線資訊

連線資料庫的 DNS 設置在 .env 文件中，以保護敏感資訊。

# 會員系統

會員系統建置中...

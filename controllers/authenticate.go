package controllers

import (
	"context"
	"http-server/config"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// 定義一個專用的上下文鍵類型
type contextKey string

// 定義特定的鍵
const UsernameContextKey contextKey = "username"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "session-name") // 獲取 Session
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			// 未登入，重定向到登入頁面
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 將用戶名存入 Context，供後續處理使用
		ctx := context.WithValue(r.Context(), UsernameContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

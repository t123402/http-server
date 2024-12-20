package applications

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"http-server/config"
	"http-server/models"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 驗證輸入
	if req.Username == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// 如果 Role 是空值，給它一個預設值 1
	if req.Role == "" {
		req.Role = "1"
	}

	// 加密密碼
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	// 新增用戶
	err = models.AddUser(req.Username, string(hashedPassword), req.Role)
	if err != nil {
		if sql.ErrNoRows == err {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 查詢用戶
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// 驗證密碼
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 保存到 Session
	session, _ := config.Store.Get(r, "session-name") // 創建/獲取 Session
	session.Values["username"] = req.Username  // 保存用戶名到 Session
	session.Values["id"] = "1"
	session.Values["nickname"] = "小可"
	session.Values["gender"] = "M"
	seserr := session.Save(r, w)                  // 保存 Session
	if seserr != nil {
		fmt.Printf("Session save error: %v\n", seserr)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// 返回成功響應
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session-name")
	session.Options.MaxAge = -1 // 設置過期時間，刪除 Session
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 定義 MeHandler，返回用戶 Session 資訊
func MeHandler(w http.ResponseWriter, r *http.Request) {
	// 從 Context 中取得用戶名
	username, ok := r.Context().Value(UsernameContextKey).(string)
	if !ok || username == "" {
		http.Error(w, "未登入", http.StatusUnauthorized)
		return
	}

	// 假設 Session 中保存了其他用戶資訊
	session, _ := config.Store.Get(r, "session-name")
	userID, _ := session.Values["id"].(string)
	nickname, _ := session.Values["nickname"].(string)
	gender, _ := session.Values["gender"].(string)

	// 返回用戶資訊
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"id":       userID,
		"username": username,
		"nickname": nickname,
		"gender":   gender,
	})
}
package applications

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"http-server/config"
	"http-server/models"
	"net/http"
	"strings"
	"time"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
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

	// 解析生日
	var birthday *time.Time // 使用指針處理非必填情況
	if req.Birthday != "" {
		parsedBirthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			fmt.Printf("Birthday parsing error: %v\n", err)
			http.Error(w, "Invalid birthday format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		birthday = &parsedBirthday
	}

	// 驗證輸入
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	// 加密密碼
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	// 新增用戶 Role給它一個預設值 87
	err = models.AddUser(req.Username, string(hashedPassword), "87")
	if err != nil {
		if sql.ErrNoRows == err {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}

	// 新增用戶資訊
	err = models.AddProfile(req.Username, req.Nickname, req.Firstname, req.Lastname, req.Email, req.Gender, birthday)
	if err != nil {
		if sql.ErrNoRows == err {
			http.Error(w, "Profile already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create profile", http.StatusInternalServerError)
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
			http.Error(w, "User Database error", http.StatusInternalServerError)
		}
		return
	}

	// 驗證密碼
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 查詢用戶資訊
	profile, err := models.GetProfileByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Profile not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Profile Database error", http.StatusInternalServerError)
		}
		return
	}

	// 查詢用戶角色
	role, err := models.GetRoleById(user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Role not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Role Database error", http.StatusInternalServerError)
		}
		return
	}

	// 保存到 Session
	session, _ := config.Store.Get(r, "session-name") // 創建/獲取 Session
	session.Values["username"] = user.Username        // 保存用戶名到 Session
	session.Values["id"] = user.ID
	session.Values["nickname"] = profile.Nickname
	session.Values["roleid"] = role.ID
	session.Values["rolename"] = role.Name
	session.Values["gender"] = profile.Gender
	seserr := session.Save(r, w) // 保存 Session
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
	userID, _ := session.Values["id"].(int)
	nickname, _ := session.Values["nickname"].(string)
	roleid, _ := session.Values["roleid"].(int)
	rolename, _ := session.Values["rolename"].(string)
	gender, _ := session.Values["gender"].(string)

	// 返回用戶資訊
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"id":       fmt.Sprintf("%d", userID),
		"username": username,
		"nickname": nickname,
		"roleid":   fmt.Sprintf("%d", roleid),
		"rolename": rolename,
		"gender":   gender,
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/auth/profile/")
	if username == "" {
		// 如果未提供 Username，返回 HTTP 400 錯誤
		http.Error(w, "Missing Username", http.StatusBadRequest)
		return
	}

	// 查詢用戶資訊
	profile, err := models.GetProfileByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Profile not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Profile Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Nickname  string `json:"nickname"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"`
	}{
		Nickname:  profile.Nickname,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		Email:     profile.Email,
		Gender:    profile.Gender,
		Birthday:  profile.Birthday.Format("2006-01-02"), // 格式化日期
	})
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
}

// 修改資料
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 確認請求方法是否為 PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/auth/profile/update/")
	if username == "" {
		// 如果未提供 Username，返回 HTTP 400 錯誤
		http.Error(w, "Missing Username", http.StatusBadRequest)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 解析生日
	var birthday *time.Time
	if req.Birthday != "" {
		parsedBirthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			fmt.Printf("Birthday parsing error: %v\n", err)
			http.Error(w, "Invalid birthday format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		birthday = &parsedBirthday
	}

	// 更新用戶資訊
	if err := models.UpdateProfileByUsername(username, req.Nickname, req.Firstname, req.Lastname, req.Email, req.Gender, birthday); err != nil {
		// 更新失敗，返回 HTTP 500 錯誤
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// 獲取當前 Session
	session, _ := config.Store.Get(r, "session-name")
	session.Values["nickname"] = req.Nickname
	session.Values["gender"] = req.Gender
	seserr := session.Save(r, w) // 保存 Session
	if seserr != nil {
		fmt.Printf("Session save error: %v\n", seserr)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// 返回 HTTP 200 OK，表示更新成功
	w.WriteHeader(http.StatusOK)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/auth/change-password/")
	if username == "" {
		// 如果未提供 Username，返回 HTTP 400 錯誤
		http.Error(w, "Missing Username", http.StatusBadRequest)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 查詢用戶
	user, err := models.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "User Database error", http.StatusInternalServerError)
		}
		return
	}

	// 驗證密碼
	if !CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 加密密碼
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	// 更新用戶密碼
	if err := models.ChangePasswordByUsername(username, hashedPassword); err != nil {
		// 更新失敗，返回 HTTP 500 錯誤
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	// 返回成功響應
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Change password successful")
}

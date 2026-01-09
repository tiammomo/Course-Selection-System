package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID string `json:"user_id"`
}

// WhoAmIResponse 获取当前用户响应
type WhoAmIResponse struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	UserType int    `json:"user_type"`
}

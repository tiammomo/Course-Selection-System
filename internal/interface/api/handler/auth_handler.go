package handler

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"course_select/internal/application/dto"
	"course_select/internal/domain/model"
	"course_select/internal/domain/service"
	"course_select/internal/pkg/errcode"
	"course_select/internal/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
	sessionKey  string
	cookieName  string
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService, sessionKey, cookieName string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		sessionKey:  sessionKey,
		cookieName:  cookieName,
	}
}

// Login 登录
// @Summary 用户登录
// @Description 用户使用用户名密码登录
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录请求"
// @Success 200 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	member, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	// 创建 Session
	session := sessions.Default(c)
	sessionID := h.authService.GenerateSessionID()
	h.authService.CreateSession(session, member)
	if err := session.Save(); err != nil {
		c.JSON(200, response.Fail(errcode.UnknownError.WithMsg("会话保存失败")))
		return
	}

	// 设置 Cookie
	c.SetCookie(h.cookieName, sessionID, 3600, "/", "", false, true)

	c.JSON(200, response.Success(map[string]string{
		"user_id": strconv.Itoa(member.UserID),
	}))
}

// Logout 登出
// @Summary 用户登出
// @Description 用户退出登录
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie(h.cookieName)
	if err != nil {
		c.JSON(200, response.Fail(errcode.LoginRequired))
		return
	}

	session := sessions.Default(c)

	// 检查 Session 是否存在
	if session.Get(sessionID) == nil {
		c.JSON(200, response.Fail(errcode.LoginRequired))
		return
	}

	session.Delete(sessionID)
	if err := session.Save(); err != nil {
		c.JSON(200, response.Fail(errcode.UnknownError.WithMsg("会话保存失败")))
		return
	}

	// 删除 Cookie
	c.SetCookie(h.cookieName, sessionID, -1, "/", "", false, true)

	c.JSON(200, response.Success(nil))
}

// WhoAmI 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的信息
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response
// @Router /auth/whoami [get]
func (h *AuthHandler) WhoAmI(c *gin.Context) {
	sessionData, exists := c.Get("session_data")
	if !exists {
		c.JSON(200, response.Fail(errcode.LoginRequired))
		return
	}

	data, ok := sessionData.(map[string]interface{})
	if !ok {
		c.JSON(200, response.Fail(errcode.UnknownError.WithMsg("无效的会话数据")))
		return
	}

	userIDVal, ok := data["user_id"]
	if !ok {
		c.JSON(200, response.Fail(errcode.UnknownError.WithMsg("无效的用户ID")))
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(200, response.Fail(errcode.UnknownError.WithMsg("无效的用户ID")))
		return
	}

	nickname := ""
	nicknameVal, ok := data["nickname"]
	if ok {
		if s, ok := nicknameVal.(string); ok {
			nickname = s
		}
	}
	username := ""
	usernameVal, ok := data["username"]
	if ok {
		if s, ok := usernameVal.(string); ok {
			username = s
		}
	}
	userType := 0
	userTypeVal, ok := data["user_type"]
	if ok {
		if t, ok := userTypeVal.(int); ok {
			userType = t
		}
	}

	c.JSON(200, response.Success(dto.WhoAmIResponse{
		UserID:   userID,
		Nickname: nickname,
		Username: username,
		UserType: userType,
	}))
}

// GetUserIDFromSession 从 Session 获取用户ID
func (h *AuthHandler) GetUserIDFromSession(c *gin.Context) (string, bool) {
	sessionData, exists := c.Get("session_data")
	if !exists {
		return "", false
	}

	data, ok := sessionData.(map[string]interface{})
	if !ok {
		return "", false
	}

	userIDVal, ok := data["user_id"]
	if !ok {
		return "", false
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return "", false
	}
	return userID, userID != ""
}

// GetUserTypeFromSession 从 Session 获取用户类型
func (h *AuthHandler) GetUserTypeFromSession(c *gin.Context) (model.UserType, bool) {
	sessionData, exists := c.Get("session_data")
	if !exists {
		return 0, false
	}

	data, ok := sessionData.(map[string]interface{})
	if !ok {
		return 0, false
	}

	userTypeVal, ok := data["user_type"]
	if !ok {
		return 0, false
	}
	userType, ok := userTypeVal.(int)
	if !ok {
		return 0, false
	}
	return model.UserType(userType), true
}

// GetUserIDFromSession 从 Gin Context 获取用户ID (包级别辅助函数)
func GetUserIDFromSession(c *gin.Context) (string, bool) {
	sessionData, exists := c.Get("session_data")
	if !exists {
		return "", false
	}

	data, ok := sessionData.(map[string]interface{})
	if !ok {
		return "", false
	}

	userIDVal, ok := data["user_id"]
	if !ok {
		return "", false
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return "", false
	}
	return userID, userID != ""
}

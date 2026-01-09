package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"course_select/internal/domain/model"
	"course_select/internal/pkg/response"
	"course_select/internal/domain/service"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService *service.AuthService
	sessionKey  string
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService *service.AuthService, sessionKey string) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		sessionKey:  sessionKey,
	}
}

// RequireAuth 需要认证
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie("camp-session")
		if err != nil {
			c.JSON(200, response.Unauthorized("用户未登录"))
			c.Abort()
			return
		}

		session := sessions.Default(c)
		v := session.Get(m.sessionKey)
		if v == nil {
			c.JSON(200, response.Unauthorized("用户未登录"))
			c.Abort()
			return
		}

		c.Set("session_id", sessionId)
		c.Set("session_data", v)
		c.Next()
	}
}

// RequireAdmin 需要管理员权限
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, exists := c.Get("session_data")
		if !exists {
			c.JSON(200, response.Unauthorized("用户未登录"))
			c.Abort()
			return
		}

		sessionData, ok := data.(map[string]interface{})
		if !ok {
			c.JSON(200, response.Unauthorized("无效的会话"))
			c.Abort()
			return
		}

		userTypeVal, ok := sessionData["user_type"]
		if !ok {
			c.JSON(200, response.Unauthorized("无效的会话"))
			c.Abort()
			return
		}
		userType, ok := userTypeVal.(int)
		if !ok || userType != int(model.UserTypeAdmin) {
			c.JSON(200, response.Forbidden("没有操作权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireStudent 需要学生权限
func (m *AuthMiddleware) RequireStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, exists := c.Get("session_data")
		if !exists {
			c.JSON(200, response.Unauthorized("用户未登录"))
			c.Abort()
			return
		}

		sessionData, ok := data.(map[string]interface{})
		if !ok {
			c.JSON(200, response.Unauthorized("无效的会话"))
			c.Abort()
			return
		}

		userTypeVal, ok := sessionData["user_type"]
		if !ok {
			c.JSON(200, response.Unauthorized("无效的会话"))
			c.Abort()
			return
		}
		userType, ok := userTypeVal.(int)
		if !ok || userType != int(model.UserTypeStudent) {
			c.JSON(200, response.Forbidden("需要学生权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"course_select/internal/pkg/logger"
)

// LoggerMiddleware 日志中间件
type LoggerMiddleware struct{}

// NewLoggerMiddleware 创建日志中间件
func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

// ZapLogger 日志中间件
func (m *LoggerMiddleware) ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("HTTP Request",
			logger.Int("status", status),
			logger.String("method", method),
			logger.String("path", path),
			logger.Duration("latency", latency),
			logger.String("client_ip", c.ClientIP()),
		)
	}
}

// Recovery 恢复中间件
func (m *LoggerMiddleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", logger.Any("error", err))
				c.AbortWithStatusJSON(500, gin.H{
					"code": 500,
					"msg":  "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}

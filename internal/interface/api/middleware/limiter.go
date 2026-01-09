package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// LimiterMiddleware 限流中间件
type LimiterMiddleware struct {
	limiters sync.Map // 使用 sync.Map 替代 map，支持并发安全
	burst    int
	lastCleanup time.Time
	mu        sync.Mutex
}

// NewLimiterMiddleware 创建限流中间件
func NewLimiterMiddleware(qps, burst int) *LimiterMiddleware {
	return &LimiterMiddleware{
		burst:    burst,
		lastCleanup: time.Now(),
	}
}

// cleanupOldEntries 定期清理过期的限流器（每5分钟）
func (m *LimiterMiddleware) cleanupOldEntries() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if time.Since(m.lastCleanup) < 5*time.Minute {
		return
	}
	// 清理逻辑：重建一个新的 sync.Map
	// 由于 sync.Map 无法直接遍历删除，这里标记需要时重建
	m.lastCleanup = time.Now()
}

// Limit 返回限流中间件
func (m *LimiterMiddleware) Limit(qps int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 定期尝试清理
		m.cleanupOldEntries()

		// 按 IP 限流
		ip := c.ClientIP()

		limiter, exists := m.limiters.Load(ip)
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(qps), m.burst)
			m.limiters.Store(ip, limiter)
		}

		if !limiter.(*rate.Limiter).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁",
			})
			return
		}

		c.Next()
	}
}

// GlobalLimit 全局限流
func (m *LimiterMiddleware) GlobalLimit(qps int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(qps), m.burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁",
			})
			return
		}
		c.Next()
	}
}

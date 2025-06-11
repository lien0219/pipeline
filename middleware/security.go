package middleware

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit 速率限制配置
const (
	MaxRequestsPerMinute = 100         // 每分钟最大请求数
	WindowDuration       = time.Minute // 时间窗口（1分钟）
)

// ipRequest 记录单个IP的请求计数和窗口开始时间
type ipRequest struct {
	count       int
	windowStart time.Time
}

// rateLimiter 速率限制器（并发安全）
var (
	rateLimiter = struct {
		sync.RWMutex
		ips map[string]*ipRequest
	}{
		ips: make(map[string]*ipRequest),
	}
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; font-src 'self' data:")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self'")
		}
		c.Next()
	}
}

func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录审计日志
		if c.Request.Method != http.MethodGet {
			userID := c.GetUint("userID")
			log := model.AuditLog{
				UserID:     userID,
				Action:     c.Request.Method,
				ObjectType: c.FullPath(),
				RequestIP:  c.ClientIP(),
				Details:    c.Request.URL.String(),
			}
			global.DB.Create(&log)
		}
		c.Next()
	}
}

// RateLimitMiddleware 速率限制中间件（限制同一IP每分钟最多100次请求）
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c.Request)

		rateLimiter.Lock()
		defer rateLimiter.Unlock()

		now := time.Now()
		ir, exists := rateLimiter.ips[ip]

		// 新IP或时间窗口已过期，重置计数
		if !exists || now.Sub(ir.windowStart) > WindowDuration {
			rateLimiter.ips[ip] = &ipRequest{
				count:       1,
				windowStart: now,
			}
		} else {
			ir.count++
			if ir.count > MaxRequestsPerMinute {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"code":    http.StatusTooManyRequests,
					"message": "请求过于频繁，请1分钟后再试",
				})
				return
			}
		}

		c.Next()
	}
}

// 辅助函数getClientIP 获取客户端真实IP（考虑反向代理场景）
func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	// 从RemoteAddr解析IP（格式为IP:端口）
	addr := r.RemoteAddr
	idx := strings.LastIndex(addr, ":")
	if idx > 0 {
		return addr[:idx]
	}
	return addr
}

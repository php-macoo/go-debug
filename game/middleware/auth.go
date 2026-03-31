// Package middleware 提供 Gin 中间件，当前包含 Bearer Token 认证。
package middleware

import (
	"net/http"
	"strings"

	"go-debug/game/pkg/resp"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

// ContextKeyUID 是认证中间件写入 gin.Context 的用户 ID 键名。
const ContextKeyUID = "uid"

// Auth 返回一个 Gin 中间件，从 Authorization 头提取 Bearer Token 并验证。
// 验证通过后将用户 ID 写入 Context，失败则返回 401 并终止请求链。
func Auth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			resp.Fail(c, http.StatusUnauthorized, -1, "未登录")
			c.Abort()
			return
		}
		uid, err := authSvc.ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			resp.Fail(c, http.StatusUnauthorized, -1, "登录已过期")
			c.Abort()
			return
		}
		c.Set(ContextKeyUID, uid)
		c.Next()
	}
}

// OptionalAuth 与 Auth 相同解析 Bearer Token，但不强制：无令牌或无效时直接放行（不写 uid）。
// 用于「登录优先，否则可走游客」的接口（如开局、交成绩）。
func OptionalAuth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.Next()
			return
		}
		uid, err := authSvc.ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.Next()
			return
		}
		c.Set(ContextKeyUID, uid)
		c.Next()
	}
}

// GetUID 从 gin.Context 中获取已认证的用户 ID（需配合 Auth 中间件使用）。
// 类型异常或未设置时返回 0，避免 panic。
func GetUID(c *gin.Context) int64 {
	uid, exists := c.Get(ContextKeyUID)
	if !exists {
		return 0
	}
	if id, ok := uid.(int64); ok {
		return id
	}
	return 0
}

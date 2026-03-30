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

// GetUID 从 gin.Context 中获取已认证的用户 ID（需配合 Auth 中间件使用）。
func GetUID(c *gin.Context) int64 {
	uid, _ := c.Get(ContextKeyUID)
	return uid.(int64)
}

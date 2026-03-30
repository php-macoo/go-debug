// Package resp 提供统一的 JSON 响应工具函数，确保所有接口返回格式一致。
//
// 响应格式: {"code": 0, "msg": "ok", "data": {...}}
//   - code=0 表示成功，非 0 表示业务错误
//   - msg 为人可读的消息
//   - data 仅在成功时携带
package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// R 是所有 API 响应的统一结构体。
type R struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// OK 返回 200 成功响应。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, R{Code: 0, Msg: "ok", Data: data})
}

// Fail 返回指定 HTTP 状态码和业务错误码的失败响应。
func Fail(c *gin.Context, httpCode, bizCode int, msg string) {
	c.JSON(httpCode, R{Code: bizCode, Msg: msg})
}

// Fail400 返回 400 Bad Request。
func Fail400(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, -1, msg)
}

// Fail401 返回 401 Unauthorized。
func Fail401(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, -1, msg)
}

// Fail500 返回 500 Internal Server Error。
func Fail500(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, -1, msg)
}

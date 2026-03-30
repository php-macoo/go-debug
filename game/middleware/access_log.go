package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"go-debug/game/dao"
	"go-debug/game/model"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

const maxBodyLog = 2000

// AccessLog 记录每次 API 请求的方法、路径、请求体、响应体、用户信息与耗时，异步写入数据库。
func AccessLog(logDAO *dao.ApiLogDAO) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var reqBody string
		contentType := c.GetHeader("Content-Type")
		if c.Request.Body != nil && !strings.HasPrefix(contentType, "multipart/") {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			reqBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		} else if strings.HasPrefix(contentType, "multipart/") {
			reqBody = "[multipart/form-data]"
		}

		rbw := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = rbw

		c.Next()

		var uid int64
		if v, exists := c.Get(ContextKeyUID); exists {
			if id, ok := v.(int64); ok {
				uid = id
			}
		}

		if len(reqBody) > maxBodyLog {
			reqBody = reqBody[:maxBodyLog]
		}
		respBody := rbw.body.String()
		if len(respBody) > maxBodyLog {
			respBody = respBody[:maxBodyLog]
		}

		record := &model.ApiLog{
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Query:      c.Request.URL.RawQuery,
			ReqBody:    reqBody,
			StatusCode: c.Writer.Status(),
			RespBody:   respBody,
			UserID:     uid,
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			LatencyMs:  time.Since(start).Milliseconds(),
		}
		go logDAO.Create(record)
	}
}

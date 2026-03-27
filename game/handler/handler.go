package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-debug/game/service"
)

type Handler struct {
	User  *service.UserService
	Score *service.ScoreService
	Auth  *service.AuthService
}

func New(user *service.UserService, score *service.ScoreService, auth *service.AuthService) *Handler {
	return &Handler{User: user, Score: score, Auth: auth}
}

// RegisterRoutes 将所有 API 路由挂载到 mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/register", h.Register)
	mux.HandleFunc("/api/login", h.Login)
	mux.HandleFunc("/api/user", h.GetUser)
	mux.HandleFunc("/api/score", h.SubmitScore)
	mux.HandleFunc("/api/leaderboard", h.Leaderboard)
}

// ────── 认证提取 ──────

func (h *Handler) authUID(r *http.Request) (int64, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0, fmt.Errorf("no token")
	}
	return h.Auth.ValidateToken(strings.TrimPrefix(auth, "Bearer "))
}

// ────── JSON 响应 ──────

type apiResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func writeOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResp{Code: 0, Msg: "ok", Data: data})
}

func writeFail(w http.ResponseWriter, httpCode, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(apiResp{Code: code, Msg: msg})
}

// ────── 工具 ──────

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func maskPhone(p string) string {
	if len(p) >= 11 {
		return p[:3] + "****" + p[len(p)-4:]
	}
	return p
}

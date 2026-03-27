package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) SubmitScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeFail(w, 405, -1, "method not allowed")
		return
	}
	uid, err := h.authUID(r)
	if err != nil {
		writeFail(w, 401, -1, "未登录")
		return
	}
	var req struct {
		TimeMs int `json:"completionTimeMs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TimeMs <= 0 {
		writeFail(w, 400, -1, "无效的完成时间")
		return
	}

	rank, err := h.Score.Submit(uid, req.TimeMs, r.UserAgent(), clientIP(r))
	if err != nil {
		writeFail(w, 500, -1, "提交失败")
		return
	}
	writeOK(w, map[string]interface{}{"rank": rank})
}

func (h *Handler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeFail(w, 405, -1, "method not allowed")
		return
	}
	list, err := h.Score.Leaderboard()
	if err != nil {
		writeFail(w, 500, -1, "查询失败")
		return
	}
	writeOK(w, map[string]interface{}{"list": list})
}

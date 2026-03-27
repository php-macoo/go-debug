package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeFail(w, 405, -1, "method not allowed")
		return
	}
	var req struct {
		Phone string `json:"phone"`
		Pwd   string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeFail(w, 400, -1, "请求格式错误")
		return
	}
	if len(req.Phone) < 11 {
		writeFail(w, 400, -1, "手机号格式不正确")
		return
	}
	if len(req.Pwd) < 6 {
		writeFail(w, 400, -1, "密码至少 6 位")
		return
	}

	user, token, err := h.User.Register(req.Phone, req.Pwd)
	if err != nil {
		writeFail(w, 400, -1, err.Error())
		return
	}
	writeOK(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id": user.ID, "username": user.Username, "phone": maskPhone(user.Phone),
		},
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeFail(w, 405, -1, "method not allowed")
		return
	}
	var req struct {
		Phone string `json:"phone"`
		Pwd   string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeFail(w, 400, -1, "请求格式错误")
		return
	}

	user, token, err := h.User.Login(req.Phone, req.Pwd)
	if err != nil {
		writeFail(w, 400, -1, err.Error())
		return
	}
	writeOK(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id": user.ID, "username": user.Username, "phone": maskPhone(user.Phone),
		},
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeFail(w, 405, -1, "method not allowed")
		return
	}
	uid, err := h.authUID(r)
	if err != nil {
		writeFail(w, 401, -1, "未登录")
		return
	}
	user, err := h.User.GetByID(uid)
	if err != nil {
		writeFail(w, 500, -1, "系统错误")
		return
	}
	writeOK(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id": user.ID, "username": user.Username, "phone": maskPhone(user.Phone),
		},
	})
}

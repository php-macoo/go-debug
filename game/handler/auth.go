package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go-debug/game/middleware"
	"go-debug/game/pkg/resp"
	"go-debug/game/pkg/util"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userSvc   *service.UserService
	staticDir string
}

func NewAuthHandler(userSvc *service.UserService, staticDir string) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, staticDir: staticDir}
}

// Register 处理用户注册请求，source 字段记录注册来源游戏。
// POST /api/auth/register  Body: {"phone":"...","password":"...","source":"match3"}
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Phone  string `json:"phone" binding:"required"`
		Pwd    string `json:"password" binding:"required"`
		Source string `json:"source"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "请求格式错误")
		return
	}
	if len(req.Phone) < 11 {
		resp.Fail400(c, "手机号格式不正确")
		return
	}
	if len(req.Pwd) < 6 {
		resp.Fail400(c, "密码至少 6 位")
		return
	}

	user, token, err := h.userSvc.Register(req.Phone, req.Pwd, req.Source)
	if err != nil {
		resp.Fail400(c, err.Error())
		return
	}
	resp.OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"username": user.Username,
			"phone":  util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// Login 处理用户登录请求。
// POST /api/auth/login  Body: {"phone":"...","password":"..."}
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Pwd   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "请求格式错误")
		return
	}

	user, token, err := h.userSvc.Login(req.Phone, req.Pwd)
	if err != nil {
		resp.Fail400(c, err.Error())
		return
	}
	resp.OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"username": user.Username,
			"phone":  util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// GetUser 获取当前登录用户信息。
// GET /api/user
func (h *AuthHandler) GetUser(c *gin.Context) {
	uid := middleware.GetUID(c)
	user, err := h.userSvc.GetByID(uid)
	if err != nil {
		resp.Fail500(c, "系统错误")
		return
	}
	resp.OK(c, gin.H{
		"user": gin.H{
			"username": user.Username,
			"phone":  util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// UpdateProfile 更新用户名。
// PUT /api/user/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	uid := middleware.GetUID(c)
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "请求格式错误")
		return
	}
	if len(req.Username) < 1 || len(req.Username) > 20 {
		resp.Fail400(c, "用户名长度 1-20 个字符")
		return
	}

	user, err := h.userSvc.UpdateUsername(uid, req.Username)
	if err != nil {
		resp.Fail500(c, err.Error())
		return
	}
	resp.OK(c, gin.H{
		"user": gin.H{
			"username": user.Username,
			"phone":  util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// UploadAvatar 上传头像图片。
// POST /api/user/avatar
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	uid := middleware.GetUID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		resp.Fail400(c, "请选择图片文件")
		return
	}
	defer file.Close()

	if header.Size > 2*1024*1024 {
		resp.Fail400(c, "图片大小不能超过 2MB")
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".png"
	}

	avatarDir := filepath.Join(h.staticDir, "avatars")
	os.MkdirAll(avatarDir, 0o755)

	filename := fmt.Sprintf("%d_%d%s", uid, time.Now().UnixMilli(), ext)
	savePath := filepath.Join(avatarDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		resp.Fail500(c, "保存文件失败")
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		resp.Fail500(c, "保存文件失败")
		return
	}

	avatarURL := "/avatars/" + filename
	user, err := h.userSvc.UpdateAvatar(uid, avatarURL)
	if err != nil {
		resp.Fail500(c, err.Error())
		return
	}
	resp.OK(c, gin.H{
		"user": gin.H{
			"username": user.Username,
			"phone":  util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// SetAvatar 设置预置头像（emoji）。
// PUT /api/user/avatar
func (h *AuthHandler) SetAvatar(c *gin.Context) {
	uid := middleware.GetUID(c)
	var req struct {
		Avatar string `json:"avatar" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "请求格式错误")
		return
	}

	user, err := h.userSvc.UpdateAvatar(uid, req.Avatar)
	if err != nil {
		resp.Fail500(c, err.Error())
		return
	}
	resp.OK(c, gin.H{
		"user": gin.H{
			"username": user.Username,
			"phone":  util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

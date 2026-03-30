// Package handler 是 HTTP 处理层，每个方法只负责：参数绑定 → 调用 Service → 返回响应。
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

// AuthHandler 处理用户认证和个人资料相关的 HTTP 请求。
type AuthHandler struct {
	userSvc   *service.UserService
	staticDir string
}

// NewAuthHandler 创建 AuthHandler，注入 UserService 和静态文件目录路径（用于头像存储）。
func NewAuthHandler(userSvc *service.UserService, staticDir string) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, staticDir: staticDir}
}

// Register 处理用户注册请求。
// POST /api/register  Body: {"phone":"13800001111","password":"123456"}
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Pwd   string `json:"password" binding:"required"`
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

	user, token, err := h.userSvc.Register(req.Phone, req.Pwd)
	if err != nil {
		resp.Fail400(c, err.Error())
		return
	}
	resp.OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"id": user.ID, "username": user.Username,
			"phone": util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// Login 处理用户登录请求。
// POST /api/login  Body: {"phone":"13800001111","password":"123456"}
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
			"id": user.ID, "username": user.Username,
			"phone": util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// GetUser 获取当前登录用户的信息（需认证中间件）。
// GET /api/user  Header: Authorization: Bearer <token>
func (h *AuthHandler) GetUser(c *gin.Context) {
	uid := middleware.GetUID(c)
	user, err := h.userSvc.GetByID(uid)
	if err != nil {
		resp.Fail500(c, "系统错误")
		return
	}
	resp.OK(c, gin.H{
		"user": gin.H{
			"id": user.ID, "username": user.Username,
			"phone": util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// UpdateProfile 更新用户名。
// PUT /api/user/profile  Body: {"username":"newName"}
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
			"id": user.ID, "username": user.Username,
			"phone": util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// UploadAvatar 上传头像图片，保存到 static/avatars/ 目录。
// POST /api/user/avatar  multipart/form-data  field: file
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
			"id": user.ID, "username": user.Username,
			"phone": util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

// SetAvatar 设置预置头像（emoji）。
// PUT /api/user/avatar  Body: {"avatar":"🐱"}
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
			"id": user.ID, "username": user.Username,
			"phone": util.MaskPhone(user.Phone), "avatar": user.Avatar,
		},
	})
}

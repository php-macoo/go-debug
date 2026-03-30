package handler

import (
	"go-debug/game/middleware"
	"go-debug/game/pkg/resp"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

// ScoreHandler 处理游戏成绩相关的 HTTP 请求（提交成绩/排行榜）。
type ScoreHandler struct {
	scoreSvc *service.ScoreService
}

// NewScoreHandler 创建 ScoreHandler，注入 ScoreService。
func NewScoreHandler(scoreSvc *service.ScoreService) *ScoreHandler {
	return &ScoreHandler{scoreSvc: scoreSvc}
}

// Submit 提交游戏成绩（需认证中间件）。
// POST /api/score  Body: {"completionTimeMs":12345}
func (h *ScoreHandler) Submit(c *gin.Context) {
	uid := middleware.GetUID(c)
	var req struct {
		TimeMs int `json:"completionTimeMs" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "无效的完成时间")
		return
	}

	rank, err := h.scoreSvc.Submit(uid, req.TimeMs, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		resp.Fail500(c, "提交失败")
		return
	}
	resp.OK(c, gin.H{"rank": rank})
}

// Leaderboard 获取 Top 10 排行榜（需认证中间件）。
// GET /api/leaderboard  Header: Authorization: Bearer <token>
func (h *ScoreHandler) Leaderboard(c *gin.Context) {
	list, err := h.scoreSvc.Leaderboard()
	if err != nil {
		resp.Fail500(c, "查询失败")
		return
	}
	resp.OK(c, gin.H{"list": list})
}

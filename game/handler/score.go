package handler

import (
	"go-debug/game/middleware"
	"go-debug/game/pkg/resp"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

type ScoreHandler struct {
	scoreSvc *service.ScoreService
}

func NewScoreHandler(scoreSvc *service.ScoreService) *ScoreHandler {
	return &ScoreHandler{scoreSvc: scoreSvc}
}

// Submit 提交游戏成绩，gameKey 从 URL 路径参数获取。
// POST /api/game/:gameKey/score
func (h *ScoreHandler) Submit(c *gin.Context) {
	uid := middleware.GetUID(c)
	gameKey := c.Param("gameKey")
	var req struct {
		TimeMs int `json:"completionTimeMs" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "无效的完成时间")
		return
	}

	rank, err := h.scoreSvc.Submit(uid, gameKey, req.TimeMs, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		resp.Fail500(c, "提交失败")
		return
	}
	resp.OK(c, gin.H{"rank": rank})
}

// Leaderboard 获取指定游戏的 Top 10 排行榜。
// GET /api/game/:gameKey/leaderboard
func (h *ScoreHandler) Leaderboard(c *gin.Context) {
	gameKey := c.Param("gameKey")
	list, err := h.scoreSvc.Leaderboard(gameKey)
	if err != nil {
		resp.Fail500(c, "查询失败")
		return
	}
	resp.OK(c, gin.H{"list": list})
}

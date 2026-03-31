package handler

import (
	"errors"
	"strings"

	"go-debug/game/middleware"
	"go-debug/game/pkg/resp"
	"go-debug/game/service"

	"github.com/gin-gonic/gin"
)

type ScoreHandler struct {
	scoreSvc *service.ScoreService
	userSvc  *service.UserService
}

func NewScoreHandler(scoreSvc *service.ScoreService, userSvc *service.UserService) *ScoreHandler {
	return &ScoreHandler{scoreSvc: scoreSvc, userSvc: userSvc}
}

// playUserID 已登录用 token 中的 uid；否则要求请求头 X-Guest-Device-Id（客户端本地持久化的匿名 ID）。
func (h *ScoreHandler) playUserID(c *gin.Context) (int64, error) {
	if uid := middleware.GetUID(c); uid != 0 {
		return uid, nil
	}
	key := strings.TrimSpace(c.GetHeader("X-Guest-Device-Id"))
	if key == "" {
		return 0, errPlayNeedIdentity
	}
	return h.userSvc.EnsureGuest(key)
}

// errPlayNeedIdentity 表示既无有效登录也无游客头。
var errPlayNeedIdentity = errors.New("need auth or guest device id")

// StartRun 创建对局凭证，用于后续成绩上报的服务端计时。
// POST /api/game/:gameKey/run/start
func (h *ScoreHandler) StartRun(c *gin.Context) {
	uid, err := h.playUserID(c)
	if err != nil {
		if errors.Is(err, errPlayNeedIdentity) {
			resp.Fail400(c, "请登录，或由客户端在请求头携带 X-Guest-Device-Id 以匿名游玩")
			return
		}
		resp.Fail400(c, err.Error())
		return
	}
	gameKey := c.Param("gameKey")
	runID, err := h.scoreSvc.StartRun(uid, gameKey)
	if err != nil {
		var se *service.ScoreSubmitError
		if errors.As(err, &se) {
			resp.Fail(c, se.Status, -1, se.Msg)
			return
		}
		resp.Fail500(c, "创建对局失败")
		return
	}
	resp.OK(c, gin.H{"runId": runID})
}

// Submit 提交游戏成绩，需携带先调用 run/start 获得的 runId。
// POST /api/game/:gameKey/score
func (h *ScoreHandler) Submit(c *gin.Context) {
	uid, err := h.playUserID(c)
	if err != nil {
		if errors.Is(err, errPlayNeedIdentity) {
			resp.Fail400(c, "请登录，或由客户端在请求头携带 X-Guest-Device-Id 以匿名游玩")
			return
		}
		resp.Fail400(c, err.Error())
		return
	}
	gameKey := c.Param("gameKey")
	var req struct {
		RunID  string `json:"runId" binding:"required"`
		TimeMs int    `json:"completionTimeMs" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail400(c, "请求格式错误：需要 runId 与 completionTimeMs")
		return
	}

	rank, err := h.scoreSvc.Submit(uid, gameKey, req.RunID, req.TimeMs, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		var se *service.ScoreSubmitError
		if errors.As(err, &se) {
			resp.Fail(c, se.Status, -1, se.Msg)
			return
		}
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

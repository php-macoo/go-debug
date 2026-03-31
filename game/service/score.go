package service

import (
	"errors"
	"net/http"
	"time"

	"go-debug/game/config"
	"go-debug/game/dao"
	"go-debug/game/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ScoreSubmitError 业务可预期的上报失败，携带 HTTP 状态码供 handler 返回。
type ScoreSubmitError struct {
	Status int
	Msg    string
}

func (e *ScoreSubmitError) Error() string { return e.Msg }

type ScoreService struct {
	db       *gorm.DB
	scoreDAO *dao.ScoreDAO
	gameDAO  *dao.GameDAO
	runDAO   *dao.GameRunDAO
	cfg      *config.Config
}

func NewScoreService(db *gorm.DB, scoreDAO *dao.ScoreDAO, gameDAO *dao.GameDAO, runDAO *dao.GameRunDAO, cfg *config.Config) *ScoreService {
	return &ScoreService{db: db, scoreDAO: scoreDAO, gameDAO: gameDAO, runDAO: runDAO, cfg: cfg}
}

// StartRun 创建一局对局凭证（需先校验游戏已上线）。
func (s *ScoreService) StartRun(userID int64, gameKey string) (runID string, err error) {
	if _, err := s.gameDAO.GetOnlineByKey(gameKey); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", &ScoreSubmitError{http.StatusNotFound, "游戏不存在或未上线"}
		}
		return "", err
	}
	ttl := time.Duration(s.cfg.Score.RunTTLMinutes) * time.Minute
	return s.runDAO.Create(userID, gameKey, ttl)
}

// Submit 校验 gameKey、对局凭证、上报间隔与完成时间上下限后写入成绩；完成时间取 min(客户端上报, 服务端墙钟) 再夹在配置范围内。
func (s *ScoreService) Submit(userID int64, gameKey, runID string, clientReportedMs int, ua, ip string) (rank int, err error) {
	if runID == "" {
		return 0, &ScoreSubmitError{http.StatusBadRequest, "缺少对局凭证 runId"}
	}
	if _, err := s.gameDAO.GetOnlineByKey(gameKey); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, &ScoreSubmitError{http.StatusNotFound, "游戏不存在或未上线"}
		}
		return 0, err
	}
	minMs, maxMs := s.cfg.CompletionLimits(gameKey)
	if clientReportedMs <= 0 {
		return 0, &ScoreSubmitError{http.StatusBadRequest, "无效的完成时间"}
	}
	if clientReportedMs < minMs {
		return 0, &ScoreSubmitError{http.StatusBadRequest, "完成时间过短"}
	}
	if maxMs > 0 && clientReportedMs > maxMs {
		return 0, &ScoreSubmitError{http.StatusBadRequest, "完成时间过长"}
	}

	interval := time.Duration(s.cfg.Score.SubmitMinIntervalSeconds) * time.Second
	var finalMs int

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var run model.GameRun
		q := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("run_id = ? AND user_id = ? AND game_key = ? AND used_at IS NULL AND expires_at > ?",
				runID, userID, gameKey, time.Now())
		if err := q.First(&run).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ScoreSubmitError{http.StatusBadRequest, "对局凭证无效、已使用或已过期"}
			}
			return err
		}

		var recent int64
		if err := tx.Model(&model.Score{}).
			Where("user_id = ? AND game_key = ? AND created_at > ?", userID, gameKey, time.Now().Add(-interval)).
			Count(&recent).Error; err != nil {
			return err
		}
		if recent > 0 {
			return &ScoreSubmitError{http.StatusTooManyRequests, "提交过于频繁，请稍后再试"}
		}

		serverMs := int(time.Since(run.StartedAt).Milliseconds())
		if serverMs < 0 {
			serverMs = 0
		}
		candidate := clientReportedMs
		if serverMs < candidate {
			candidate = serverMs
		}
		finalMs = candidate
		if finalMs < minMs {
			finalMs = minMs
		}
		if maxMs > 0 && finalMs > maxMs {
			finalMs = maxMs
		}

		now := time.Now()
		run.UsedAt = &now
		if err := tx.Save(&run).Error; err != nil {
			return err
		}
		score := &model.Score{
			UserID:           userID,
			GameKey:          gameKey,
			CompletionTimeMs: finalMs,
			UserAgent:        ua,
			IP:               ip,
		}
		return s.scoreDAO.CreateWithTx(tx, score)
	})
	if err != nil {
		return 0, err
	}

	r, err := s.scoreDAO.GetRank(gameKey, finalMs)
	if err != nil {
		return 1, nil
	}
	return r, nil
}

// Leaderboard 获取指定游戏的 Top 10 排行榜。
func (s *ScoreService) Leaderboard(gameKey string) ([]model.LeaderboardEntry, error) {
	return s.scoreDAO.TopN(gameKey, 10)
}

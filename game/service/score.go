package service

import (
	"go-debug/game/dao"
	"go-debug/game/model"
)

// ScoreService 封装成绩提交和排行榜查询的业务逻辑。
type ScoreService struct {
	scoreDAO *dao.ScoreDAO
}

// NewScoreService 创建 ScoreService，注入 ScoreDAO。
func NewScoreService(scoreDAO *dao.ScoreDAO) *ScoreService {
	return &ScoreService{scoreDAO: scoreDAO}
}

// Submit 提交游戏成绩并返回该成绩在排行榜中的排名。
func (s *ScoreService) Submit(userID int64, timeMs int, ua, ip string) (int, error) {
	score := &model.Score{
		UserID:           userID,
		CompletionTimeMs: timeMs,
		UserAgent:        ua,
		IP:               ip,
	}
	if err := s.scoreDAO.Create(score); err != nil {
		return 0, err
	}
	rank, err := s.scoreDAO.GetRank(timeMs)
	if err != nil {
		return 1, nil
	}
	return rank, nil
}

// Leaderboard 获取 Top 10 排行榜数据。
func (s *ScoreService) Leaderboard() ([]model.LeaderboardEntry, error) {
	return s.scoreDAO.TopN(10)
}

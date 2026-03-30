package service

import (
	"go-debug/game/dao"
	"go-debug/game/model"
)

type ScoreService struct {
	scoreDAO *dao.ScoreDAO
}

func NewScoreService(scoreDAO *dao.ScoreDAO) *ScoreService {
	return &ScoreService{scoreDAO: scoreDAO}
}

// Submit 提交游戏成绩并返回排名，gameKey 标识具体游戏。
func (s *ScoreService) Submit(userID int64, gameKey string, timeMs int, ua, ip string) (int, error) {
	score := &model.Score{
		UserID:           userID,
		GameKey:          gameKey,
		CompletionTimeMs: timeMs,
		UserAgent:        ua,
		IP:               ip,
	}
	if err := s.scoreDAO.Create(score); err != nil {
		return 0, err
	}
	rank, err := s.scoreDAO.GetRank(gameKey, timeMs)
	if err != nil {
		return 1, nil
	}
	return rank, nil
}

// Leaderboard 获取指定游戏的 Top 10 排行榜。
func (s *ScoreService) Leaderboard(gameKey string) ([]model.LeaderboardEntry, error) {
	return s.scoreDAO.TopN(gameKey, 10)
}

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

func (s *ScoreService) Leaderboard() ([]model.LeaderboardEntry, error) {
	return s.scoreDAO.TopN(10)
}

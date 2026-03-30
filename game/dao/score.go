package dao

import (
	"go-debug/game/model"

	"gorm.io/gorm"
)

type ScoreDAO struct {
	db *gorm.DB
}

func NewScoreDAO(db *gorm.DB) *ScoreDAO {
	return &ScoreDAO{db: db}
}

func (d *ScoreDAO) Create(score *model.Score) error {
	return d.db.Create(score).Error
}

// GetRank 计算给定完成时间在指定游戏中所有玩家最佳成绩中的排名。
func (d *ScoreDAO) GetRank(gameKey string, timeMs int) (int, error) {
	var rank int64
	err := d.db.Raw(`
		SELECT COUNT(*)+1 FROM (
			SELECT MIN(completion_time_ms) AS best FROM scores WHERE game_key = ? GROUP BY user_id
		) t WHERE t.best < ?`, gameKey, timeMs).Scan(&rank).Error
	if err != nil {
		return 0, err
	}
	return int(rank), nil
}

// TopN 获取指定游戏的排行榜前 n 名。
func (d *ScoreDAO) TopN(gameKey string, n int) ([]model.LeaderboardEntry, error) {
	var list []model.LeaderboardEntry
	err := d.db.Raw(`
		SELECT u.username, MIN(s.completion_time_ms) AS best_ms
		FROM scores s JOIN users u ON s.user_id = u.id
		WHERE s.game_key = ?
		GROUP BY s.user_id, u.username
		ORDER BY best_ms ASC
		LIMIT ?`, gameKey, n).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].Rank = i + 1
	}
	return list, nil
}

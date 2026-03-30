package dao

import (
	"go-debug/game/model"

	"gorm.io/gorm"
)

// ScoreDAO 封装成绩表的数据库操作。
type ScoreDAO struct {
	db *gorm.DB
}

// NewScoreDAO 创建 ScoreDAO 实例。
func NewScoreDAO(db *gorm.DB) *ScoreDAO {
	return &ScoreDAO{db: db}
}

// Create 插入一条成绩记录。
func (d *ScoreDAO) Create(score *model.Score) error {
	return d.db.Create(score).Error
}

// GetRank 计算给定完成时间在所有玩家最佳成绩中的排名（1-based）。
func (d *ScoreDAO) GetRank(timeMs int) (int, error) {
	var rank int64
	err := d.db.Raw(`
		SELECT COUNT(*)+1 FROM (
			SELECT MIN(completion_time_ms) AS best FROM scores GROUP BY user_id
		) t WHERE t.best < ?`, timeMs).Scan(&rank).Error
	if err != nil {
		return 0, err
	}
	return int(rank), nil
}

// TopN 获取排行榜前 n 名，按每位玩家最佳成绩升序排列。
func (d *ScoreDAO) TopN(n int) ([]model.LeaderboardEntry, error) {
	var list []model.LeaderboardEntry
	err := d.db.Raw(`
		SELECT u.username, MIN(s.completion_time_ms) AS best_ms
		FROM scores s JOIN users u ON s.user_id = u.id
		GROUP BY s.user_id, u.username
		ORDER BY best_ms ASC
		LIMIT ?`, n).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].Rank = i + 1
	}
	return list, nil
}

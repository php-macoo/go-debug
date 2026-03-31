package dao

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"go-debug/game/model"

	"gorm.io/gorm"
)

type GameRunDAO struct {
	db *gorm.DB
}

func NewGameRunDAO(db *gorm.DB) *GameRunDAO {
	return &GameRunDAO{db: db}
}

func randomRunID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Create 为用户对某游戏新建一局，返回 runId（随机十六进制串）。
func (d *GameRunDAO) Create(userID int64, gameKey string, ttl time.Duration) (runID string, err error) {
	runID, err = randomRunID()
	if err != nil {
		return "", err
	}
	now := time.Now()
	run := model.GameRun{
		RunID:     runID,
		UserID:    userID,
		GameKey:   gameKey,
		StartedAt: now,
		ExpiresAt: now.Add(ttl),
	}
	return runID, d.db.Create(&run).Error
}

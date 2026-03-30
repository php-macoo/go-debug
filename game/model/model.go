// Package model 定义数据库表对应的 GORM 模型及 API 用到的数据传输对象。
package model

import "time"

// User 对应 users 表，存储玩家账号信息。
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string    `gorm:"uniqueIndex;type:varchar(20);not null" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	Avatar    string    `gorm:"type:varchar(500);default:''" json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Score 对应 scores 表，记录每次游戏完成的成绩和客户端信息。
type Score struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"index;not null" json:"userId"`
	CompletionTimeMs int       `gorm:"column:completion_time_ms;index;not null" json:"completionTimeMs"`
	UserAgent        string    `gorm:"type:text" json:"userAgent"`
	IP               string    `gorm:"type:varchar(45)" json:"ip"`
	CreatedAt        time.Time `json:"createdAt"`
	User             User      `gorm:"foreignKey:UserID" json:"-"`
}

// LeaderboardEntry 是排行榜单条记录的 DTO，不直接映射数据库表。
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	BestMs   int    `json:"bestTimeMs"`
}

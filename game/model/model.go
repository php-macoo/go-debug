package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string    `gorm:"uniqueIndex;type:varchar(20);not null" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Score struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"index;not null" json:"userId"`
	CompletionTimeMs int       `gorm:"column:completion_time_ms;index;not null" json:"completionTimeMs"`
	UserAgent        string    `gorm:"type:text" json:"userAgent"`
	IP               string    `gorm:"type:varchar(45)" json:"ip"`
	CreatedAt        time.Time `json:"createdAt"`
	User             User      `gorm:"foreignKey:UserID" json:"-"`
}

type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	BestMs   int    `json:"bestTimeMs"`
}

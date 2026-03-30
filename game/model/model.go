package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string    `gorm:"uniqueIndex;type:varchar(20);not null" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	Avatar    string    `gorm:"type:varchar(500);default:''" json:"avatar"`
	Source    string    `gorm:"type:varchar(50);default:''" json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Score struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"index;not null" json:"userId"`
	GameKey          string    `gorm:"type:varchar(50);index;not null;default:'match3'" json:"gameKey"`
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

type Game struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	GameKey   string    `gorm:"uniqueIndex;type:varchar(50);not null" json:"gameKey"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Icon      string    `gorm:"type:varchar(50)" json:"icon"`
	Desc      string    `gorm:"type:varchar(500)" json:"desc"`
	URL       string    `gorm:"type:varchar(200)" json:"url"`
	Status    string    `gorm:"type:varchar(20);not null;default:'online'" json:"status"`
	SortOrder int       `gorm:"default:0" json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ApiLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Method     string    `gorm:"type:varchar(10);not null" json:"method"`
	Path       string    `gorm:"type:varchar(500);not null" json:"path"`
	Query      string    `gorm:"type:text" json:"query"`
	ReqBody    string    `gorm:"type:text" json:"reqBody"`
	StatusCode int       `gorm:"not null" json:"statusCode"`
	RespBody   string    `gorm:"type:text" json:"respBody"`
	UserID     int64     `gorm:"index;default:0" json:"userId"`
	IP         string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent  string    `gorm:"type:text" json:"userAgent"`
	LatencyMs  int64     `gorm:"not null" json:"latencyMs"`
	CreatedAt  time.Time `json:"createdAt"`
}

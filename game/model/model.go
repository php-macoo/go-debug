package model

import "time"

// User 用户表。
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键" json:"id"`
	Phone     string    `gorm:"uniqueIndex;type:varchar(20);not null;comment:手机号" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null;comment:密码哈希" json:"-"`
	Username  string    `gorm:"type:varchar(50);not null;comment:昵称" json:"username"`
	Avatar    string    `gorm:"type:varchar(500);default:'';comment:头像URL或emoji" json:"avatar"`
	Source    string    `gorm:"type:varchar(50);default:'';comment:注册来源游戏标识" json:"source"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updatedAt"`
}

// Score 用户游戏成绩记录表。
type Score struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键" json:"id"`
	UserID           int64     `gorm:"index;not null;type:bigint unsigned;comment:用户ID" json:"userId"`
	GameKey          string    `gorm:"type:varchar(50);index;not null;default:'match3';comment:游戏标识" json:"gameKey"`
	CompletionTimeMs int       `gorm:"column:completion_time_ms;index;not null;comment:完成用时毫秒" json:"completionTimeMs"`
	UserAgent        string    `gorm:"type:text;comment:上报时User-Agent" json:"userAgent"`
	IP               string    `gorm:"type:varchar(45);comment:上报时客户端IP" json:"ip"`
	CreatedAt        time.Time `gorm:"comment:记录创建时间" json:"createdAt"`
}

// LeaderboardEntry 排行榜查询结果，非数据库表。
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	BestMs   int    `json:"bestTimeMs"`
}

// Game 游戏大厅配置表。
type Game struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键" json:"id"`
	GameKey   string    `gorm:"uniqueIndex;type:varchar(50);not null;comment:游戏唯一标识" json:"gameKey"`
	Name      string    `gorm:"type:varchar(100);not null;comment:展示名称" json:"name"`
	Icon      string    `gorm:"type:varchar(50);comment:封面图标或emoji" json:"icon"`
	Desc      string    `gorm:"type:varchar(500);comment:简介" json:"desc"`
	URL       string    `gorm:"type:varchar(200);comment:入口相对路径" json:"url"`
	Status    string    `gorm:"type:varchar(20);not null;default:'online';comment:状态 online/coming/offline" json:"status"`
	SortOrder int       `gorm:"default:0;comment:排序权重升序" json:"sortOrder"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updatedAt"`
}

// GameRun 对局凭证表，用于成绩上报的服务端计时与一次性校验。
type GameRun struct {
	ID        int64      `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键" json:"-"`
	RunID     string     `gorm:"uniqueIndex;type:char(36);not null;comment:对局凭证ID" json:"runId"`
	UserID    int64      `gorm:"index:idx_game_run_user_game,priority:1;not null;type:bigint unsigned;comment:用户ID" json:"-"`
	GameKey   string     `gorm:"type:varchar(50);index:idx_game_run_user_game,priority:2;not null;comment:游戏标识" json:"-"`
	StartedAt time.Time  `gorm:"not null;comment:对局开始时间" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;comment:凭证过期时间" json:"-"`
	UsedAt    *time.Time `gorm:"comment:成绩已提交占用时间" json:"-"`
	CreatedAt time.Time  `gorm:"comment:记录创建时间" json:"-"`
}

// ApiLog API 访问日志表。
type ApiLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键" json:"id"`
	Method     string    `gorm:"type:varchar(10);not null;comment:HTTP方法" json:"method"`
	Path       string    `gorm:"type:varchar(500);not null;comment:请求路径" json:"path"`
	Query      string    `gorm:"type:text;comment:查询串" json:"query"`
	ReqBody    string    `gorm:"type:text;comment:请求体截断" json:"reqBody"`
	StatusCode int       `gorm:"not null;comment:响应状态码" json:"statusCode"`
	RespBody   string    `gorm:"type:text;comment:响应体截断" json:"respBody"`
	UserID     int64     `gorm:"index;default:0;type:bigint unsigned;comment:已登录用户ID未登录为0" json:"userId"`
	IP         string    `gorm:"type:varchar(45);comment:客户端IP" json:"ip"`
	UserAgent  string    `gorm:"type:text;comment:User-Agent" json:"userAgent"`
	LatencyMs  int64     `gorm:"not null;comment:处理耗时毫秒" json:"latencyMs"`
	CreatedAt  time.Time `gorm:"comment:记录时间" json:"createdAt"`
}

// Package config 定义应用配置结构体，并从 YAML 文件加载配置。
package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 是应用的顶层配置，包含服务器、数据库和认证三部分。
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	Score    ScoreConfig    `yaml:"score"`
}

// ScoreConfig 成绩上报：限流、对局凭证、各游戏完成时间上下限。
type ScoreConfig struct {
	SubmitMinIntervalSeconds int `yaml:"submit_min_interval_seconds"`
	RunTTLMinutes            int `yaml:"run_ttl_minutes"`
	// 未在 games 中单独配置时使用
	DefaultMinCompletionTimeMs int `yaml:"default_min_completion_time_ms"`
	DefaultMaxCompletionTimeMs int `yaml:"default_max_completion_time_ms"` // 0 表示不限制上限
	Games                      map[string]GameScoreLimit `yaml:"games"`
}

// GameScoreLimit 单个 gameKey 的完成时间毫秒上下限；0 表示沿用 default。
type GameScoreLimit struct {
	MinCompletionTimeMs int `yaml:"min_completion_time_ms"`
	MaxCompletionTimeMs int `yaml:"max_completion_time_ms"`
}

// ServerConfig 定义 HTTP 服务相关配置。
type ServerConfig struct {
	Addr      string `yaml:"addr"`
	StaticDir string `yaml:"static_dir"`
}

// DatabaseConfig 定义 MySQL 数据库连接参数。
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	// LogSQL 为 true 时向标准输出打印 GORM 执行的 SQL（开发环境用；生产务必关闭）。
	LogSQL bool `yaml:"log_sql"`
}

// DSN 生成 MySQL 连接字符串。withDB=true 时包含数据库名，false 时不包含（用于建库）。
func (d DatabaseConfig) DSN(withDB bool) string {
	name := ""
	if withDB {
		name = d.Name
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		d.User, d.Password, d.Host, d.Port, name)
}

// AuthConfig 定义认证相关配置（token 密钥和有效期）。
type AuthConfig struct {
	TokenSecret     string `yaml:"token_secret"`
	TokenExpireDays int    `yaml:"token_expire_days"`
}

// MustLoad 从指定路径加载 YAML 配置文件，失败时直接 Fatal 退出。
func MustLoad(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("读取配置文件 %s 失败: %v", path, err)
	}
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	cfg.normalizeScore()
	return &cfg
}

func (c *Config) normalizeScore() {
	if c.Score.SubmitMinIntervalSeconds <= 0 {
		c.Score.SubmitMinIntervalSeconds = 15
	}
	if c.Score.RunTTLMinutes <= 0 {
		c.Score.RunTTLMinutes = 120
	}
	if c.Score.DefaultMinCompletionTimeMs <= 0 {
		c.Score.DefaultMinCompletionTimeMs = 5000
	}
	if c.Score.Games == nil {
		c.Score.Games = make(map[string]GameScoreLimit)
	}
}

// CompletionLimits 返回某 gameKey 生效的 [minMs, maxMs]；maxMs==0 表示无上限。
func (c *Config) CompletionLimits(gameKey string) (minMs int, maxMs int) {
	minMs = c.Score.DefaultMinCompletionTimeMs
	maxMs = c.Score.DefaultMaxCompletionTimeMs
	if g, ok := c.Score.Games[gameKey]; ok {
		if g.MinCompletionTimeMs > 0 {
			minMs = g.MinCompletionTimeMs
		}
		if g.MaxCompletionTimeMs > 0 {
			maxMs = g.MaxCompletionTimeMs
		}
	}
	return minMs, maxMs
}

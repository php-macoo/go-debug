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
	return &cfg
}

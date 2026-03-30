// Package dao 是数据访问层，封装所有数据库操作（GORM）。
package dao

import (
	"database/sql"
	"fmt"

	"go-debug/game/config"
	"go-debug/game/model"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接：先用原生连接确保数据库存在，再通过 GORM 建立连接并自动迁移表结构。
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// 先用不带数据库名的 DSN 连接，确保目标数据库存在
	rawDB, err := sql.Open("mysql", cfg.DSN(false))
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}
	_, _ = rawDB.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.Name)
	rawDB.Close()

	// 使用 GORM 连接目标数据库
	db, err := gorm.Open(mysql.Open(cfg.DSN(true)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", cfg.Name, err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// 自动迁移：根据模型定义创建/更新表结构
	if err = db.AutoMigrate(&model.User{}, &model.Score{}); err != nil {
		return nil, fmt.Errorf("迁移失败: %w", err)
	}
	return db, nil
}

// CloseDB 安全关闭数据库连接。
func CloseDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

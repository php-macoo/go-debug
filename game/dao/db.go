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

func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	rawDB, err := sql.Open("mysql", cfg.DSN(false))
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}
	_, _ = rawDB.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.Name)
	rawDB.Close()

	db, err := gorm.Open(mysql.Open(cfg.DSN(true)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", cfg.Name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = db.AutoMigrate(&model.User{}, &model.Score{}, &model.Game{}, &model.ApiLog{}); err != nil {
		return nil, fmt.Errorf("迁移失败: %w", err)
	}
	return db, nil
}

func CloseDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

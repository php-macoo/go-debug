package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"go-debug/game/config"

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

	var gormLog logger.Interface
	if cfg.LogSQL {
		gormLog = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	} else {
		gormLog = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN(true)), &gorm.Config{
		Logger: gormLog,
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

	return db, nil
}

func CloseDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

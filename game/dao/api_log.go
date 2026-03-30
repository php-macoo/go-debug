package dao

import (
	"go-debug/game/model"

	"gorm.io/gorm"
)

type ApiLogDAO struct {
	db *gorm.DB
}

func NewApiLogDAO(db *gorm.DB) *ApiLogDAO {
	return &ApiLogDAO{db: db}
}

func (d *ApiLogDAO) Create(record *model.ApiLog) error {
	return d.db.Create(record).Error
}

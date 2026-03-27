package dao

import (
	"go-debug/game/model"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) Create(user *model.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) FindByPhone(phone string) (*model.User, error) {
	var u model.User
	if err := d.db.Where("phone = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *UserDAO) FindByID(id int64) (*model.User, error) {
	var u model.User
	if err := d.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

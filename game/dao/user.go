package dao

import (
	"go-debug/game/model"

	"gorm.io/gorm"
)

// UserDAO 封装用户表的数据库操作。
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建 UserDAO 实例。
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// Create 插入一条用户记录，成功后 user.ID 会被回填。
func (d *UserDAO) Create(user *model.User) error {
	return d.db.Create(user).Error
}

// FindByPhone 根据手机号查询用户，未找到时返回 gorm.ErrRecordNotFound。
func (d *UserDAO) FindByPhone(phone string) (*model.User, error) {
	var u model.User
	if err := d.db.Where("phone = ?", phone).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID 根据主键 ID 查询用户。
func (d *UserDAO) FindByID(id int64) (*model.User, error) {
	var u model.User
	if err := d.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateAvatar 更新用户头像路径。
func (d *UserDAO) UpdateAvatar(id int64, avatar string) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Update("avatar", avatar).Error
}

// UpdateUsername 更新用户名。
func (d *UserDAO) UpdateUsername(id int64, username string) error {
	return d.db.Model(&model.User{}).Where("id = ?", id).Update("username", username).Error
}

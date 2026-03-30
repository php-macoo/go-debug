package service

import (
	"errors"
	"fmt"
	"strings"

	"go-debug/game/dao"
	"go-debug/game/model"

	"gorm.io/gorm"
)

// UserService 封装用户注册、登录、查询、资料修改等业务逻辑。
type UserService struct {
	userDAO *dao.UserDAO
	auth    *AuthService
}

// NewUserService 创建 UserService，注入 DAO 和认证服务。
func NewUserService(userDAO *dao.UserDAO, auth *AuthService) *UserService {
	return &UserService{userDAO: userDAO, auth: auth}
}

// Register 注册新用户：哈希密码 → 写入数据库 → 签发 token。
// 默认用户名为 "user" + 手机号后四位。
func (s *UserService) Register(phone, password string) (*model.User, string, error) {
	hashed, err := s.auth.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("系统错误")
	}

	user := &model.User{
		Phone:    phone,
		Password: hashed,
		Username: "user" + phone[len(phone)-4:],
	}
	if err = s.userDAO.Create(user); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return nil, "", fmt.Errorf("该手机号已注册")
		}
		return nil, "", fmt.Errorf("注册失败")
	}

	token := s.auth.GenerateToken(user.ID)
	return user, token, nil
}

// Login 登录：校验手机号 → 比对密码 → 签发 token。
func (s *UserService) Login(phone, password string) (*model.User, string, error) {
	user, err := s.userDAO.FindByPhone(phone)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", fmt.Errorf("手机号未注册")
	}
	if err != nil {
		return nil, "", fmt.Errorf("系统错误")
	}
	if s.auth.CheckPassword(user.Password, password) != nil {
		return nil, "", fmt.Errorf("密码错误")
	}

	token := s.auth.GenerateToken(user.ID)
	return user, token, nil
}

// GetByID 根据用户 ID 查询用户信息。
func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.userDAO.FindByID(id)
}

// UpdateAvatar 更新用户头像并返回最新用户信息。
func (s *UserService) UpdateAvatar(id int64, avatar string) (*model.User, error) {
	if err := s.userDAO.UpdateAvatar(id, avatar); err != nil {
		return nil, fmt.Errorf("更新头像失败")
	}
	return s.userDAO.FindByID(id)
}

// UpdateUsername 更新用户名并返回最新用户信息。
func (s *UserService) UpdateUsername(id int64, username string) (*model.User, error) {
	if err := s.userDAO.UpdateUsername(id, username); err != nil {
		return nil, fmt.Errorf("更新用户名失败")
	}
	return s.userDAO.FindByID(id)
}

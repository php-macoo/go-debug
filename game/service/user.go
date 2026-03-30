package service

import (
	"errors"
	"fmt"
	"strings"

	"go-debug/game/dao"
	"go-debug/game/model"

	"gorm.io/gorm"
)

type UserService struct {
	userDAO *dao.UserDAO
	auth    *AuthService
}

func NewUserService(userDAO *dao.UserDAO, auth *AuthService) *UserService {
	return &UserService{userDAO: userDAO, auth: auth}
}

// Register 注册新用户，source 记录注册来源（游戏编号）。
func (s *UserService) Register(phone, password, source string) (*model.User, string, error) {
	hashed, err := s.auth.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("系统错误")
	}

	user := &model.User{
		Phone:    phone,
		Password: hashed,
		Username: "user" + phone[len(phone)-4:],
		Source:   source,
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

func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.userDAO.FindByID(id)
}

func (s *UserService) UpdateAvatar(id int64, avatar string) (*model.User, error) {
	if err := s.userDAO.UpdateAvatar(id, avatar); err != nil {
		return nil, fmt.Errorf("更新头像失败")
	}
	return s.userDAO.FindByID(id)
}

func (s *UserService) UpdateUsername(id int64, username string) (*model.User, error) {
	if err := s.userDAO.UpdateUsername(id, username); err != nil {
		return nil, fmt.Errorf("更新用户名失败")
	}
	return s.userDAO.FindByID(id)
}

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

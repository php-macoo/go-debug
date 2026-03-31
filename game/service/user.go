package service

import (
	"crypto/sha256"
	"encoding/hex"
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

// guestSyntheticPhone 将客户端持有的设备标识映射为占位的唯一 phone（≤20 字符，满足 users.phone 长度）。
func guestSyntheticPhone(deviceKey string) string {
	sum := sha256.Sum256([]byte(deviceKey))
	hexStr := hex.EncodeToString(sum[:])
	return "g" + hexStr[:19]
}

// EnsureGuest 根据匿名设备标识查找或创建一条用户记录（Source=guest），用于未登录完局与上榜。
func (s *UserService) EnsureGuest(deviceKey string) (int64, error) {
	deviceKey = strings.TrimSpace(deviceKey)
	if len(deviceKey) < 8 || len(deviceKey) > 128 {
		return 0, fmt.Errorf("游客标识长度须在 8～128 之间")
	}
	phone := guestSyntheticPhone(deviceKey)
	user, err := s.userDAO.FindByPhone(phone)
	if err == nil {
		return user.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	hashed, err := s.auth.HashPassword(deviceKey + "|guest")
	if err != nil {
		return 0, fmt.Errorf("系统错误")
	}
	u := &model.User{
		Phone:    phone,
		Password: hashed,
		Username: "游客",
		Source:   "guest",
	}
	if err := s.userDAO.Create(u); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			user2, err2 := s.userDAO.FindByPhone(phone)
			if err2 == nil {
				return user2.ID, nil
			}
		}
		return 0, fmt.Errorf("创建游客失败")
	}
	return u.ID, nil
}

// Package service 是业务逻辑层，编排 DAO 调用，不依赖 HTTP 层。
package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-debug/game/config"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 负责 token 签发/验证和密码哈希。
// token 格式: base64url(payload).hex(hmac-sha256)，轻量级自定义 JWT。
type AuthService struct {
	secret     string
	expireDays int
}

// NewAuthService 根据配置创建 AuthService。
func NewAuthService(cfg config.AuthConfig) *AuthService {
	return &AuthService{secret: cfg.TokenSecret, expireDays: cfg.TokenExpireDays}
}

// tokenPayload 是 token 中携带的数据。
type tokenPayload struct {
	UID int64 `json:"uid"`
	Exp int64 `json:"exp"`
}

// GenerateToken 为指定用户 ID 签发一个带有效期的 token。
func (s *AuthService) GenerateToken(uid int64) string {
	p := tokenPayload{
		UID: uid,
		Exp: time.Now().Add(time.Duration(s.expireDays) * 24 * time.Hour).Unix(),
	}
	data, _ := json.Marshal(p)
	encoded := base64.RawURLEncoding.EncodeToString(data)
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(encoded))
	sig := hex.EncodeToString(mac.Sum(nil))
	return encoded + "." + sig
}

// ValidateToken 验证 token 的签名和有效期，返回用户 ID。
func (s *AuthService) ValidateToken(token string) (int64, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("bad token format")
	}
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(parts[0]))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expected)) {
		return 0, fmt.Errorf("invalid signature")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, err
	}
	var p tokenPayload
	if err = json.Unmarshal(raw, &p); err != nil {
		return 0, err
	}
	if time.Now().Unix() > p.Exp {
		return 0, fmt.Errorf("token expired")
	}
	return p.UID, nil
}

// HashPassword 使用 bcrypt 对明文密码进行哈希。
func (s *AuthService) HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 校验明文密码是否与 bcrypt 哈希匹配。
func (s *AuthService) CheckPassword(hashed, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
}

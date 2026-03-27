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

type AuthService struct {
	secret     string
	expireDays int
}

func NewAuthService(cfg config.AuthConfig) *AuthService {
	return &AuthService{secret: cfg.TokenSecret, expireDays: cfg.TokenExpireDays}
}

type tokenPayload struct {
	UID int64 `json:"uid"`
	Exp int64 `json:"exp"`
}

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

func (s *AuthService) HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) CheckPassword(hashed, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
}

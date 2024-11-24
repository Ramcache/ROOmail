package auth

import (
	"ROOmail/pkg/db"
	"ROOmail/pkg/utils"
	"fmt"
	"sync"
)

type AuthService struct {
	blacklist sync.Map
}

var instance *AuthService
var once sync.Once

func NewAuthService() *AuthService {
	return &AuthService{}
}

func AuthServiceInstance() *AuthService {
	once.Do(func() {
		instance = &AuthService{}
	})
	return instance
}

// Проверка учетных данных пользователя
func (s *AuthService) AuthenticateUser(username, password string) (*db.User, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// Аннулирование токена (добавление в черный список)
func (s *AuthService) RevokeToken(token string) {
	s.blacklist.Store(token, true)
}

// Проверка, является ли токен аннулированным
func (s *AuthService) IsTokenRevoked(token string) bool {
	_, revoked := s.blacklist.Load(token)
	return revoked
}

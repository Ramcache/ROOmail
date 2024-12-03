package auth

import (
	"ROOmail/internal/models"
	"ROOmail/pkg/db"
	"ROOmail/pkg/utils"
	"ROOmail/pkg/utils/JWT"
	"context"
	"fmt"
	"sync"
	"time"
)

type AuthInterface interface {
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	IsTokenRevoked(token string) bool
	CleanupRevokedTokens()
	RevokeToken(ctx context.Context, token string) error
}

type AuthService struct {
	blacklist sync.Map
}

var instance *AuthService
var once sync.Once

func AuthServiceInstance() *AuthService {
	once.Do(func() {
		instance = &AuthService{}
	})
	return instance
}

func (s *AuthService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить аутентификацию пользователя")
	}

	if !utils.CheckPassword(password, user.Password) {
		time.Sleep(1 * time.Second)
		return nil, fmt.Errorf("неверные учетные данные")
	}

	return user, nil
}

func (s *AuthService) RevokeToken(ctx context.Context, token string) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context canceled or deadline exceeded")
	}

	s.blacklist.Store(token, true)
	return nil
}

func (s *AuthService) IsTokenRevoked(token string) bool {
	_, revoked := s.blacklist.Load(token)
	return revoked
}

func (s *AuthService) CleanupRevokedTokens() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C

		s.blacklist.Range(func(key, value interface{}) bool {
			token, ok := key.(string)
			if !ok {
				return true
			}

			if JWT.IsTokenExpired(token) {
				s.blacklist.Delete(token)
			}

			return true
		})
	}
}

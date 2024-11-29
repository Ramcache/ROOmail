package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ExtractUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("Отсутствует заголовок авторизации")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неверный метод подписи: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		// Используем общий подход для проверки истечения срока действия и других ошибок.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, fmt.Errorf("Токен истек")
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return 0, fmt.Errorf("Неправильный формат токена")
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, fmt.Errorf("Неверная подпись токена")
		} else {
			return 0, fmt.Errorf("Ошибка парсинга токена: %v", err)
		}
	}

	if token.Valid {
		return claims.UserID, nil
	}
	return 0, fmt.Errorf("Недействительный токен")
}

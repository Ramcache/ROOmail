package JWT

import (
	"fmt"
	"time"
)

func IsTokenExpired(token string) bool {
	// Здесь добавляется логика извлечения времени истечения токена.
	// Для примера, пусть токен хранит время истечения в секундах через разделитель ":"
	// tokenFormat: "<token>:<expiryUnixTimestamp>"
	var expiryTimestamp int64
	_, err := fmt.Sscanf(token, "%*s:%d", &expiryTimestamp)
	if err != nil {
		return false // Если не удалось извлечь время истечения, предполагаем, что токен не истек
	}

	expiryTime := time.Unix(expiryTimestamp, 0)
	return time.Now().After(expiryTime)
}

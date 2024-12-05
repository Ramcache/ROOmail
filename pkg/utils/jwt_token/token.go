package jwt_token

import (
	"fmt"
	"time"
)

func IsTokenExpired(token string) bool {
	var expiryTimestamp int64
	_, err := fmt.Sscanf(token, "%*s:%d", &expiryTimestamp)
	if err != nil {
		return false
	}

	expiryTime := time.Unix(expiryTimestamp, 0)
	return time.Now().After(expiryTime)
}

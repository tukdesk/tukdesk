package helpers

import (
	"time"

	"github.com/tukdesk/httputils/tools"
)

const (
	FileTokenType          = "file"
	FileTokenExpirationSec = 60 * 30
	FileTokenExpiration    = FileTokenExpirationSec * time.Second
)

func NewInternalFileToken(userId, key string, expiration time.Duration) string {
	data := map[string]interface{}{
		"type":   FileTokenType,
		"userId": userId,
	}

	return tools.GenerateToken(data, expiration, []byte(key))
}

func InternalFileTokenValid(userId, s, key string) bool {
	token, err := tools.ParseToken(s, []byte(key))
	if err != nil || !token.Valid {
		return false
	}

	if typeStr, _ := token.Claims["type"].(string); typeStr != FileTokenType {
		return false
	}

	if userIdStr, _ := token.Claims["userId"]; userIdStr != userId {
		return false
	}

	return true
}

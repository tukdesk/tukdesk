package helpers

import (
	"time"

	"github.com/tukdesk/httputils/tools"
)

const (
	AttachmentTokenType          = "attachment"
	AttachmentTokenExpirationSec = 60 * 30
	AttachmentTokenExpiration    = AttachmentTokenExpirationSec * time.Second
)

func NewInternalAttachmentToken(userId, key string, expiration time.Duration) string {
	data := map[string]interface{}{
		"type":   AttachmentTokenType,
		"userId": userId,
	}

	return tools.GenerateToken(data, expiration, []byte(key))
}

func InternalAttachmentTokenValid(userId, s, key string) bool {
	token, err := tools.ParseToken(s, []byte(key))
	if err != nil || !token.Valid {
		return false
	}

	if typeStr, _ := token.Claims["type"].(string); typeStr != AttachmentTokenType {
		return false
	}

	if userIdStr, _ := token.Claims["userId"]; userIdStr != userId {
		return false
	}

	return true
}

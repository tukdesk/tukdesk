package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"
)

func OutputProfile(user *models.User) M {
	return M{
		"id":       user.Id,
		"base":     user.Base,
		"personal": user.Personal,
		"created":  user.Created,
		"updated":  user.Updated,
	}
}

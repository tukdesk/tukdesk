package apis

import (
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

func AuthorizedAsAgent(user *models.User) bool {
	return user != nil && helpers.UserIsAgent(user)
}

package apis

import (
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"gopkg.in/mgo.v2/bson"
)

func AuthorizedLogged(user *models.User) bool {
	return user != nil
}

func AuthorizedAsSpecifiedUser(user *models.User, id bson.ObjectId) bool {
	return AuthorizedLogged(user) && user.Id == id
}

func AuthorizedAsAgent(user *models.User) bool {
	return AuthorizedLogged(user) && helpers.UserIsAgent(user)
}

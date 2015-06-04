package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

type OutputUser struct {
	Id       bson.ObjectId            `json:"id"`
	Base     *models.UserBaseInfo     `json:"base,omitempty"`
	Personal *models.UserPersonalInfo `json:"personal,omitempty"`
	Created  int64                    `json:"created,omitempty"`
	Updated  int64                    `json:"updated,omitempty"`
}

func OutputUserProfileInfo(user *models.User) *OutputUser {
	return &OutputUser{
		Id:       user.Id,
		Base:     &user.Base,
		Personal: &user.Personal,
		Created:  user.Created,
		Updated:  user.Updated,
	}
}

func OutputUserBaseInfo(user *models.User) *OutputUser {
	return &OutputUser{
		Id:   user.Id,
		Base: &user.Base,
	}
}

func OutputUserBaseInfoByUserId(userId interface{}) (*OutputUser, error) {
	user, err := UserFindById(userId)
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	return OutputUserBaseInfo(user), nil
}

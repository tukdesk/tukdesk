package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

type OutputUserInfo struct {
	Id       bson.ObjectId            `json:"id"`
	Base     *models.UserBaseInfo     `json:"base,omitempty"`
	Personal *models.UserPersonalInfo `json:"personal,omitempty"`
	Created  int64                    `json:"created,omitempty"`
	Updated  int64                    `json:"updated,omitempty"`
}

func OutputUserProfileInfo(user *models.User) *OutputUserInfo {
	return &OutputUserInfo{
		Id:       user.Id,
		Base:     &user.Base,
		Personal: &user.Personal,
		Created:  user.Created,
		Updated:  user.Updated,
	}
}

func OutputUserBaseInfo(user *models.User) *OutputUserInfo {
	return &OutputUserInfo{
		Id:   user.Id,
		Base: &user.Base,
	}
}

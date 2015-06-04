package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

type OutputUser struct {
	Id       bson.ObjectId            `json:"id"`
	Base     *models.UserBaseInfo     `json:"base,omitempty"`
	Personal *models.UserPersonalInfo `json:"personal,omitempty"`
	Business *models.UserBusinessInfo `json:"business,omitempty"`
	Created  int64                    `json:"created,omitempty"`
	Updated  int64                    `json:"updated,omitempty"`
}

func OutputUserBaseInfo(user *models.User) *OutputUser {
	return &OutputUser{
		Id:   user.Id,
		Base: &user.Base,
	}
}

func OutputUserProfileInfo(user *models.User) *OutputUser {
	output := OutputUserBaseInfo(user)
	output.Personal = &user.Personal
	output.Created = user.Created
	output.Updated = user.Updated
	return output
}

func OutputUserDetailInfo(user *models.User) *OutputUser {
	output := OutputUserProfileInfo(user)
	output.Business = &user.Business
	return output
}

func OutputUserBaseInfoByUserId(userId interface{}) (*OutputUser, error) {
	user, err := UserFindById(userId)
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	return OutputUserBaseInfo(user), nil
}

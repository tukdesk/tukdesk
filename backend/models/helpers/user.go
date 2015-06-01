package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"github.com/tukdesk/httputils/tools"
)

const (
	UserChannelAgent = "_AGENT"
	UserChannelEmail = "_EMAIL"

	UserNameMaxLength     = 20
	UserPasswordMinLength = 6

	UserRandNameLength = 6
)

func AgentInit(email, name, password, salt string) (*models.User, error) {
	channel := models.ChannelInfo{
		Name: UserChannelAgent,
		Id:   email,
	}

	user := models.NewUserWithChannel(channel)
	user.Personal.Email = email
	user.Base.Name = name
	user.Login.Password = tools.Encrypt(password, salt)

	if err := user.Insert(); err != nil {
		return nil, err
	}

	return user, nil
}

func AgentFind() (*models.User, error) {
	query := M{"channel.name": UserChannelAgent}
	user := &models.User{}
	return user, user.FindOne(query)
}

func UserCheckPassword(user *models.User, password, salt string) bool {
	if password == "" {
		return false
	}

	return tools.Encrypt(password, salt) == user.Login.Password
}

func UserIsAgent(user *models.User) bool {
	return user.Channel.Name == UserChannelAgent
}

func UserFindById(id interface{}) (*models.User, error) {
	user := &models.User{}
	return user, user.FindById(id)
}

func UserMustByChannel(chName, chId, email, name string) (*models.User, bool, error) {
	user, err := UserFindByChannel(chName, chId)
	// found
	if err == nil {
		return user, false, nil
	}

	// db error
	if !IsNotFound(err) {
		return nil, false, err
	}

	// not found
	channel := models.ChannelInfo{
		Name: chName,
		Id:   chId,
	}

	// init new doc
	doc := models.NewUserWithChannel(channel)
	doc.Personal.Email = email
	doc.Base.Name = name

	query := M{
		"channel.name": chName,
		"channel.id":   chId,
	}

	user = &models.User{}
	inserted, err := user.FindOrInsert(query, doc)
	return user, inserted, err
}

func UserFindByChannel(chName, chId string) (*models.User, error) {
	query := M{
		"channel.name": chName,
		"channel.id":   chId,
	}

	user := &models.User{}
	return user, user.FindOne(query)
}

func UserGetValidName(name string) string {
	_, s := tools.CutRune(name, UserNameMaxLength)
	return s
}

func UserGetValidNameFromEmail(email string) string {
	name, _ := tools.CutEmail(email)
	return UserGetValidName(name)
}

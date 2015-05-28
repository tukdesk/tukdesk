package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	UserCollectionName = "user"

	UserChannelAgent = "_AGENT"
	UserChannelEmail = "_EMAIL"
)

type User struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Updated int64         `json:"updated" bson:"updated"`
	Created int64         `json:"created" bson:"created"`

	Login    UserLoginInfo    `json:"-" bson:"login"`
	Channel  ChannelInfo      `json:"-" bson:"channel"`
	Base     UserBaseInfo     `json:"-" bson:"base"`
	Personal UserPersonalInfo `json:"-" bson:"personal"`
	Business UserBusinessInfo `json:"-" bson:"business"`
}

type UserLoginInfo struct {
	Password string `json:"-" bson:"password"`
}

type UserBaseInfo struct {
	Name   string `json:"name" bson:"name"`
	Avatar string `json:"avatar" bson:"avatar"`
}

type UserPersonalInfo struct {
	Email              string            `json:"email" bson:"email"`
	EmailCertificated  bool              `json:"emailCertificated" bson:"emailCertificated"`
	Mobile             string            `json:"mobile" bson:"mobile"`
	MobileCertificated bool              `json:"mobileCertificated" bson:"mobileCertificated"`
	Extend             map[string]string `json:"extend" bson:"extend"`
}

type UserBusinessInfo struct {
	Importance TypeUserImportance
	Tips       []*UserTip
}

type UserTip struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Content string        `json:"content" bson:"content"`
	Created int64         `json:"created" bson:"created"`
	Updated int64         `json:"updated" bson:"updated"`
}

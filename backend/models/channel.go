package models

type ChannelInfo struct {
	Name string      `json:"name" bson:"name"`
	Id   interface{} `json:"-" bson:"id"`
}

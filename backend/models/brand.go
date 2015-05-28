package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	BrandCollectionName = "brand"
)

type Brand struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Name    string        `json:"name" bson:"name"`
	APIKey  string        `json:"apiKey" bson:"apiKey"`
	Created int64         `json:"created" bson:"created"`
	Updated int64         `json:"updated" bson:"updated"`
	On      bool          `json:"on" bson:"on"`
}

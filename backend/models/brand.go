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
	APIKey  []byte        `json:"-" bson:"apiKey"`
	Salt    string        `json:"-" bson:"salt"`
	Created int64         `json:"created" bson:"created"`
	Updated int64         `json:"updated" bson:"updated"`
	On      bool          `json:"on" bson:"on"`
}

func NewBrand(name string) *Brand {
	now := Now().Unix()
	return &Brand{
		Id:      NewId(),
		Name:    name,
		Created: now,
		Updated: now,
		On:      true,
	}
}

func (this *Brand) Insert() error {
	return Insert(BrandCollectionName, this)
}

func (this *Brand) FindOne(query map[string]interface{}) error {
	return FindOne(BrandCollectionName, query, this)
}

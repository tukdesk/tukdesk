package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	BrandCollectionName = "brand"
)

type BrandBaseInfo struct {
	Name string `json:"name" bson:"name"`
	Logo string `json:"logo" bson:"logo"`
}

type BrandAutorizationInfo struct {
	APIKey string `json:"apiKey" bson:"apiKey"`
	Salt   string `json:"salt" bson:"salt"`
}

type BrandPreExtend struct {
	Ticket []*ExtendField `json:"ticket" bson:"ticket"`
	Client []*ExtendField `json:"client" bson:"client"`
}

type Brand struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Created int64         `json:"created" bson:"created"`
	Updated int64         `json:"updated" bson:"updated"`

	Base          BrandBaseInfo         `json:"base" bson:"base"`
	Authorization BrandAutorizationInfo `json:"-" bson:"authorization"`
	PreExtend     BrandPreExtend        `json:"-" bson:"preExtend"`
	On            bool                  `json:"-" bson:"on"`
}

func NewBrand(name string) *Brand {
	now := Now().Unix()
	return &Brand{
		Id:      NewId(),
		Created: now,
		Updated: now,
		Base: BrandBaseInfo{
			Name: name,
		},
		On: true,
	}
}

func (this *Brand) Insert() error {
	return Insert(BrandCollectionName, this)
}

func (this *Brand) FindAndModify(query, change map[string]interface{}) error {
	return FindAndModify(BrandCollectionName, query, change, this)
}

func (this *Brand) FindOne(query map[string]interface{}) error {
	return FindOne(BrandCollectionName, query, this)
}

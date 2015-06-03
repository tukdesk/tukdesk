package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	ExtendFieldInputTypeText     = "TEXT"
	ExtendFieldInputTypeTextarea = "TEXTAREA"
	ExtendFieldInputTypeSelect   = "SELECT"
	ExtendFieldInputTypeCheckbox = "CHECKBOX"
)

type ExtendField struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Label       string        `json:"label" bson:"label"`
	Description string        `json:"description" bson:"description"`
	Input       string        `json:"input" bson:"input"`
	Options     []string      `json:"options" bson:"options"`
	Created     int64         `json:"created" bson:"created"`
	Updated     int64         `json:"updated" bson:"updated"`
}

func NewExtendField() *ExtendField {
	now := Now().Unix()
	return &ExtendField{
		Id:      NewId(),
		Created: now,
		Updated: now,
	}
}

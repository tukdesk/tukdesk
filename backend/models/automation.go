package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	AutomationCollectionName = "automation"
)

type Automation struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Created int64         `json:"created" bson:"created"`
	Updated int64         `json:"updated" bson:"updated"`

	Priority int `json:"priority" bson:"priority"`

	Conditions map[bson.ObjectId][]*AutomationCondition `json:"conditions" bson:"conditions"`
	Actions    []*AutomationAction                      `json:"actions" bson:"actions"`
}

type AutomationCondition struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Operator string        `json:"operator" bson:"operator"`
	Field    string        `json:"field" bson:"field"`
	ValStr   string        `json:"valStr" bson:"valStr"`
}

type AutomationAction struct {
	Id     bson.ObjectId `json:"id" bson:"_id"`
	Type   string        `json:"type" bson:"type"`
	Target string        `json:"target" bson:"target"`
	ValStr string        `json:"valStr" bson:"valStr"`
}

package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	FocusCollectionName = "focus"

	FocusTypeReminder     = "REMINDER"
	FocusTypeNotification = "NOTIFICATION"
)

type Focus struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Type     string        `json:"type" bson:"type"`
	Priority TypePriority  `json:"priority" bson:"priority"`
	Status   string        `json:"status" bson:"status"`
	Resource Resource      `json:"resource" bson:"resource"`
	Message  string        `json:"message" bson:"message"`
	Created  int64         `json:"created" bson:"created"`
	Handled  int64         `json:"handled" bson:"handled"`
	Deadline int64         `json:"deadline" bson:"deadline"`
}

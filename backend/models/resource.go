package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	ResourceTypeTicket  = "TICKET"
	ResourceTypeComment = "COMMENT"
)

type Resource struct {
	Type   string        `json:"type" bson:"type"`
	Id     bson.ObjectId `json:"id" bson:"id,omitempty"`
	Parent *Resource     `json:"parent,omitempty" bson:"parent,omitempty"`
}

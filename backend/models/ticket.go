package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	TicketCollectionName = "ticket"

	TicketStatusPending     = "PENDING"
	TicketStatusReplied     = "REPLIED"
	TicketStatusResubmitted = "RESUBMITTED"
	TicketStatusDone        = "DONE"

	TicketChannelWeb   = "_WEB"
	TicketChannelEmail = "_EMAIL"
)

type Ticket struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	CreatorId bson.ObjectId `json:"-" bson:"creatorId"`
	Channel   ChannelInfo   `json:"channel" bson:"channel"`

	Priority       TypePriority      `json:"priority" bson:"priority"`
	Subject        string            `json:"subject" bson:"subject"`
	Content        string            `json:"content" bson:"content"`
	IsPublic       bool              `json:"isPublic" bson:"isPublic"`
	Created        int64             `json:"created" bson:"created"`
	Updated        int64             `json:"updated" bson:"updated"`
	FirstCommented int64             `json:"firstCommented" bson:"firstCommented"`
	Status         string            `json:"status" bson:"status"`
	Rank           int               `json:"rank" bson:"rank"`
	Extend         map[string]string `json:"extend" bson:"extend"`
}

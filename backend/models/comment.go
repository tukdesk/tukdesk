package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	CommentCollectionName = "comment"

	CommentTypeFirst    = "FIRST"
	CommentTypePublic   = "PUBLIC"
	CommentTypeFeedback = "FEEDBACK"
	CommentTypeInternal = "INTERNAL"
)

type Comment struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	TicketId    bson.ObjectId `json:"-" bson:"ticketId"`
	CreatorId   bson.ObjectId `json:"-" bson:"creatorId"`
	Type        string        `json:"type" bson:"type"`
	Content     string        `json:"content" bson:"content"`
	Created     int64         `json:"created" bson:"created"`
	Attachments []*Attachment `json:"attachmente" bson:"attachments"`
}

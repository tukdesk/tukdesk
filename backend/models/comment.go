package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	CommentCollectionName = "comment"

	CommentTypeQuestion = "QUESTION"
	CommentTypePublic   = "PUBLIC"
	CommentTypeFeedback = "FEEDBACK"
	CommentTypeInternal = "INTERNAL"
)

var (
	EmptyComment = &Comment{}
)

type Comment struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	TicketId    bson.ObjectId `json:"-" bson:"ticketId"`
	CreatorId   bson.ObjectId `json:"-" bson:"creatorId"`
	Type        string        `json:"type" bson:"type"`
	Content     string        `json:"content" bson:"content"`
	Created     int64         `json:"created" bson:"created"`
	Updated     int64         `json:"updated" bson:"updated"`
	Attachments []*File       `json:"attachmente" bson:"attachments"`
}

func NewComment(ticketId, creatorId bson.ObjectId, typ, content string) *Comment {
	now := Now().Unix()
	return &Comment{
		Id:        NewId(),
		TicketId:  ticketId,
		CreatorId: creatorId,
		Type:      typ,
		Content:   content,
		Created:   now,
		Updated:   now,
	}
}

func (this *Comment) Insert() error {
	return Insert(CommentCollectionName, this)
}

func (this *Comment) FindById(id interface{}) error {
	return FindById(CommentCollectionName, id, nil, this)
}

func (this *Comment) FindOne(query map[string]interface{}) error {
	return FindOne(CommentCollectionName, query, nil, this)
}

func (this *Comment) FindAll(query map[string]interface{}, sort []string) ([]*Comment, error) {
	all := make([]*Comment, 0, listCap)
	return all, FindAll(CommentCollectionName, query, nil, sort, &all)
}

func (this *Comment) FindAndModify(query, change map[string]interface{}) error {
	return FindAndModify(CommentCollectionName, query, nil, change, this)
}

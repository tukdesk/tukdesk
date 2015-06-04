package models

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	FocusCollectionName = "focus"

	FocusStatusPending = "PENDING"
	FocusStatusHandled = "HANDLED"
)

var (
	EmptyFocus = &Focus{}
)

type Focus struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Type     string        `json:"type" bson:"type"`
	Status   string        `json:"status" bson:"status"`
	Priority TypePriority  `json:"priority" bson:"priority"`
	Resource *Resource     `json:"resource" bson:"resource"`
	Message  string        `json:"message" bson:"message"`
	Created  int64         `json:"created" bson:"created"`
	Handled  int64         `json:"handled" bson:"handled"`
	Deadline int64         `json:"deadline" bson:"deadline"`
}

func NewFocus(typ, message string) *Focus {
	now := Now().Unix()
	return &Focus{
		Id:       NewId(),
		Type:     typ,
		Status:   FocusStatusPending,
		Priority: PriorityNormal,
		Message:  message,
		Created:  now,
	}
}

func (this *Focus) Insert() error {
	return Insert(FocusCollectionName, this)
}

func (this *Focus) FindById(id interface{}) error {
	return FindById(FocusCollectionName, id, this)
}

func (this *Focus) FindAndModify(query, change map[string]interface{}) error {
	return FindAndModify(FocusCollectionName, query, change, this)
}

func (this *Focus) Count(query map[string]interface{}) (int, error) {
	return Count(FocusCollectionName, query)
}

func (this *Focus) List(query map[string]interface{}, start, limit int, sort []string) ([]*Focus, error) {
	list := make([]*Focus, 0, listCap)
	return list, List(FocusCollectionName, query, start, limit, sort, &list)
}

func (this *Focus) UpdateById(id interface{}, change map[string]interface{}) error {
	return UpdateById(FocusCollectionName, id, change)
}

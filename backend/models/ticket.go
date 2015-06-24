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
)

var (
	EmptyTicket = &Ticket{}
)

type Ticket struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	CreatorId bson.ObjectId `json:"-" bson:"creatorId"`
	Channel   ChannelInfo   `json:"channel" bson:"channel"`

	Subject        string            `json:"subject" bson:"subject"`
	IsPublic       bool              `json:"isPublic" bson:"isPublic"`
	Created        int64             `json:"created" bson:"created"`
	Updated        int64             `json:"updated" bson:"updated"`
	FirstCommented int64             `json:"firstCommented" bson:"firstCommented"`
	Priority       TypePriority      `json:"priority" bson:"priority"`
	Status         string            `json:"status" bson:"status"`
	Rank           int               `json:"rank" bson:"rank"`
	Extend         map[string]string `json:"extend" bson:"extend"`
}

func NewTicket(creatorId bson.ObjectId) *Ticket {
	now := Now().Unix()
	return &Ticket{
		Id:        NewId(),
		CreatorId: creatorId,
		Created:   now,
		Updated:   now,
		Priority:  PriorityNormal,
		Status:    TicketStatusPending,
	}
}

func (this *Ticket) Insert() error {
	return Insert(TicketCollectionName, this)
}

func (this *Ticket) FindById(id interface{}) error {
	return FindById(TicketCollectionName, id, this)
}

func (this *Ticket) FindAndModify(query, change map[string]interface{}) error {
	return FindAndModify(TicketCollectionName, query, change, this)
}

func (this *Ticket) Count(query map[string]interface{}) (int, error) {
	return Count(TicketCollectionName, query)
}

func (this *Ticket) List(query map[string]interface{}, start, limit int, sort []string) ([]*Ticket, error) {
	list := make([]*Ticket, 0, listCap)
	return list, List(TicketCollectionName, query, start, limit, sort, &list)
}

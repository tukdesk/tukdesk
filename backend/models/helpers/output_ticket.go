package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

type OutputTicket struct {
	Id       bson.ObjectId       `json:"id"`
	Creator  *OutputUser         `json:"creator"`
	Channel  *models.ChannelInfo `json:"channel,omitempty"`
	Subject  string              `json:"subject"`
	Content  string              `json:"content"`
	Priority models.TypePriority `json:"priority,omitempty"`
	IsPublic bool                `json:"isPublic"`
	Created  int64               `json:"created"`
	Updated  int64               `json:"updated"`
	Status   string              `json:"status"`
	Rank     int                 `json:"rank,omitempty"`
	Extend   map[string]string   `json:"extend,omitempty"`
}

func OutputTicketPublicInfo(ticket *models.Ticket) (*OutputTicket, error) {
	creator, err := OutputUserBaseInfoByUserId(ticket.CreatorId)
	if err != nil {
		return nil, err
	}

	output := &OutputTicket{}
	output.Id = ticket.Id
	output.Creator = creator

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.Content = ticket.Content
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status

	return output, nil
}

func OutputTicketDetailInfo(ticket *models.Ticket) (*OutputTicket, error) {
	creator, err := OutputUserBaseInfoByUserId(ticket.CreatorId)
	if err != nil {
		return nil, err
	}

	output := &OutputTicket{}
	output.Id = ticket.Id
	output.Creator = creator

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.Content = ticket.Content
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status

	output.Priority = ticket.Priority
	output.Rank = ticket.Rank
	output.Extend = ticket.Extend

	return output, nil
}

func OutputTicketPublicInfoForList(ticket *models.Ticket) (*OutputTicket, error) {
	creator, err := OutputUserBaseInfoByUserId(ticket.CreatorId)
	if err != nil {
		return nil, err
	}

	output := &OutputTicket{}
	output.Id = ticket.Id
	output.Creator = creator

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status

	return output, nil
}

func OutputTicketDetailInfoForList(ticket *models.Ticket) (*OutputTicket, error) {
	creator, err := OutputUserBaseInfoByUserId(ticket.CreatorId)
	if err != nil {
		return nil, err
	}

	output := &OutputTicket{}
	output.Id = ticket.Id
	output.Creator = creator

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status

	output.Priority = ticket.Priority

	return output, nil
}

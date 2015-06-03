package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

type OutputTicketInfo struct {
	Id       bson.ObjectId       `json:"id"`
	Creator  *OutputUserInfo     `json:"creator"`
	Channel  *models.ChannelInfo `json:"channel,omitempty"`
	Subject  string              `json:"subject"`
	Content  string              `json:"content"`
	Priority models.TypePriority `json:"priority,omitempty"`
	IsPublic bool                `json:"isPublic"`
	Created  int64               `json:"created"`
	Updated  int64               `json:"updated"`
	Status   string              `json:"status"`
	Rank     int                 `json:"rank"`
	Extend   map[string]string   `json:"extend,omitempty"`
}

func OutputTicketPublicInfo(ticket *models.Ticket) (*OutputTicketInfo, error) {
	user, err := UserFindById(ticket.CreatorId)
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	output := &OutputTicketInfo{}
	output.Id = ticket.Id
	output.Creator = OutputUserBaseInfo(user)

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.Content = ticket.Content
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status

	return output, nil
}

func OutputTicketPublicInfoForList(ticket *models.Ticket) (*OutputTicketInfo, error) {
	user, err := UserFindById(ticket.CreatorId)
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	output := &OutputTicketInfo{}
	output.Id = ticket.Id
	output.Creator = OutputUserBaseInfo(user)

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status

	return output, nil
}

func OutputTicketDetailInfoForList(ticket *models.Ticket) (*OutputTicketInfo, error) {
	user, err := UserFindById(ticket.CreatorId)
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	output := &OutputTicketInfo{}
	output.Id = ticket.Id
	output.Creator = OutputUserBaseInfo(user)

	output.Channel = &ticket.Channel
	output.Subject = ticket.Subject
	output.Priority = ticket.Priority
	output.IsPublic = ticket.IsPublic
	output.Created = ticket.Created
	output.Updated = ticket.Updated
	output.Status = ticket.Status
	return output, nil
}

package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"github.com/tukdesk/httputils/tools"
	"gopkg.in/mgo.v2/bson"
)

const (
	TicketSubjectMaxLength     = 20
	TicketContentMinLength     = 15
	TicketExtendFieldMaxLength = 100

	TicketStatusDefault = models.TicketStatusPending

	TicketChannelWeb   = "_WEB"
	TicketChannelEmail = "_EMAIL"
)

var (
	TicketPriorityOptions = []interface{}{
		models.PriorityLow,
		models.PriorityNormal,
		models.PriorityHign,
		models.PriorityUrgent,
	}

	TicketStatusOptionsForList = []interface{}{
		models.TicketStatusPending,
		models.TicketStatusReplied,
		models.TicketStatusResubmitted,
		models.TicketStatusSolved,
	}

	TicketSortOptionsForList = []interface{}{
		"-updated",
		"updated",
		"-priority",
		"priority",
	}

	TicketStatusOptionsForCreate = []interface{}{
		models.TicketStatusPending,
		models.TicketStatusSolved,
	}

	TicketStatusOptionsForUpdate = []interface{}{
		models.TicketStatusSolved,
	}

	TicketRankOptions = []interface{}{0, 1, 2, 3, 4, 5}

	TicketChannelOptionsForCreate = []interface{}{
		TicketChannelWeb,
		TicketChannelEmail,
	}
)

func TicketGetValidSubject(s string) string {
	_, subject := tools.CutRune(s, TicketSubjectMaxLength)
	return subject
}

func TicketParseExtendFromPreSet(extend map[string]string) map[string]string {
	res := map[string]string{}

	for _, field := range currentBrand.PreExtend.Ticket {
		if val, ok := extend[field.Label]; ok {
			res[field.Label] = val
		}
	}
	return res
}

func TicketInit(ticket *models.Ticket, content string, attachments []*models.File) error {
	if err := ticket.Insert(); err != nil {
		return err
	}
	comment := models.NewComment(ticket.Id, ticket.CreatorId, models.CommentTypeQuestion, content)
	comment.Attachments = attachments
	return comment.Insert()
}

func TicketNewWithChannelName(creator *models.User, chName, subject string, extend map[string]string) *models.Ticket {
	ticket := models.NewTicket(creator.Id)
	ticket.Channel.Name = chName
	ticket.Channel.Id = ticket.Id
	ticket.Subject = subject
	ticket.Extend = extend
	return ticket
}

func TicketFindById(ticketId interface{}) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	return ticket, ticket.FindById(ticketId)
}

func TicketFindAndModify(ticket *models.Ticket, change map[string]interface{}) error {
	query := M{"_id": ticket.Id}
	return ticket.FindAndModify(query, change)
}

func TicketListAfter(query map[string]interface{}, lastId bson.ObjectId, limit int, sort []string) ([]*models.Ticket, error) {
	if !IsEmptyId(lastId) {
		if query == nil {
			query = M{}
		}
		query["_id"] = M{"$gt": lastId}
	}

	return models.EmptyTicket.List(query, 0, limit, sort)
}

func TicketCount(query map[string]interface{}) (int, error) {
	return models.EmptyTicket.Count(query)
}

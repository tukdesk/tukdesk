package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"
)

var (
	ResourceTypeOptions = []interface{}{
		models.ResourceTypeTicket,
		models.ResourceTypeComment,
	}
)

func ResourceByTicketId(ticketId interface{}) (*models.Resource, error) {
	ticket, err := TicketFindById(ticketId)
	if err != nil {
		return nil, err
	}

	return ResourceByTicket(ticket), nil
}

func ResourceByTicket(ticket *models.Ticket) *models.Resource {
	return models.NewResource(models.ResourceTypeTicket, ticket.Id, nil)
}

func ResourceByCommentId(commentId interface{}) (*models.Resource, error) {
	comment, err := CommentFindById(commentId)
	if err != nil {
		return nil, err
	}

	return ResourceByComment(comment), nil
}

func ResourceByComment(comment *models.Comment) *models.Resource {
	parent := models.NewResource(models.ResourceTypeTicket, comment.TicketId, nil)
	return models.NewResource(models.ResourceTypeComment, comment.Id, parent)
}

func ResourceByTypeAndId(typ string, id interface{}) (*models.Resource, error) {
	switch typ {
	case models.ResourceTypeTicket:
		return ResourceByTicketId(id)
	case models.ResourceTypeComment:
		return ResourceByCommentId(id)
	default:
		return nil, nil
	}
}

package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"
)

const (
	CommentContentMinLength = 5
)

var (
	CommentTypeOptionsForNonAgentView = []interface{}{
		models.CommentTypeQuestion,
		models.CommentTypePublic,
		models.CommentTypeFeedback,
	}

	CommentTypeOptionsForCreate = []interface{}{
		models.CommentTypeFeedback,
		models.CommentTypePublic,
		models.CommentTypeInternal,
	}

	CommentTypeOptionsForUpdate = []interface{}{
		models.CommentTypePublic,
	}
)

func CommentInsertForTicket(ticket *models.Ticket, creator *models.User, typ, content string) (*models.Comment, error) {
	comment := models.NewComment(ticket.Id, creator.Id, typ, content)
	return comment, comment.Insert()
}

func CommentFindByTicketIdAndCommentId(ticketId, commentId interface{}) (*models.Comment, error) {
	query := M{
		"_id":      commentId,
		"ticketId": ticketId,
	}

	comment := &models.Comment{}
	return comment, comment.FindOne(query)
}

func CommentFindAllByTicketId(ticketId interface{}, query map[string]interface{}, sort []string) ([]*models.Comment, error) {
	if query == nil {
		query = M{}
	}
	query["ticketId"] = ticketId
	return models.EmptyComment.FindAll(query, sort)
}

func CommentFindAndModify(comment *models.Comment, change map[string]interface{}) error {
	query := M{"_id": comment.Id}
	return comment.FindAndModify(query, change)
}

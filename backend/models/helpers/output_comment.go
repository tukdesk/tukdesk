package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

type OutputComment struct {
	Id      bson.ObjectId `json:"id"`
	Creator *OutputUser   `json:"creator"`
	Type    string        `json:"type"`
	Content string        `json:"content"`
	Created int64         `json:"created"`
	Updated int64         `json:"updated"`
}

func OutputCommentInfo(comment *models.Comment) (*OutputComment, error) {
	creator, err := OutputUserBaseInfoByUserId(comment.CreatorId)
	if err != nil {
		return nil, err
	}

	output := &OutputComment{}
	output.Id = comment.Id
	output.Creator = creator
	output.Type = comment.Type
	output.Content = comment.Content
	output.Created = comment.Created
	output.Updated = comment.Updated
	return output, nil
}

func OutputCommentInfos(comments []*models.Comment) ([]*OutputComment, error) {
	outputs := make([]*OutputComment, len(comments))
	if len(comments) == 0 {
		return outputs, nil
	}

	var err error

	for i, comment := range comments {
		if outputs[i], err = OutputCommentInfo(comment); err != nil {
			return nil, err
		}
	}

	return outputs, nil
}

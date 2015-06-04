package helpers

import (
	"github.com/tukdesk/tukdesk/backend/models"

	"gopkg.in/mgo.v2/bson"
)

const (
	FocusMessageMaxLength = 140

	FocusTypeReminder     = "REMINDER"
	FocusTypeNotification = "NOTIFICATION"
)

var (
	FocusPriorityOptions = []interface{}{
		models.PriorityLow,
		models.PriorityNormal,
		models.PriorityHign,
		models.PriorityUrgent,
	}

	FocusStatusOptionForList = []interface{}{
		models.FocusStatusPending,
		models.FocusStatusHandled,
	}
)

func FocusNew(typ, message string, durSec int64, priority models.TypePriority, resource *models.Resource) *models.Focus {
	focus := models.NewFocus(typ, message)
	if durSec > 0 {
		focus.Deadline = focus.Created + durSec
	}
	focus.Priority = priority
	focus.Resource = resource
	return focus
}

func FocusReminderInsert(message string, durSec int64, priority models.TypePriority, resource *models.Resource) (*models.Focus, error) {
	focus := FocusNew(FocusTypeReminder, message, durSec, priority, resource)
	return focus, focus.Insert()
}

func FocusFindById(focusId interface{}) (*models.Focus, error) {
	focus := &models.Focus{}
	return focus, focus.FindById(focusId)
}

func FocusCount(query map[string]interface{}) (int, error) {
	return models.EmptyFocus.Count(query)
}

func FocusListAfter(query map[string]interface{}, lastId bson.ObjectId, limit int, sort []string) ([]*models.Focus, error) {
	if !IsEmptyId(lastId) {
		if query == nil {
			query = M{}
		}
		query["_id"] = M{"$gt": lastId}
	}

	return models.EmptyFocus.List(query, 0, limit, sort)
}

func FocusHandled(focus *models.Focus) error {
	query := M{
		"_id": focus.Id,
	}

	change := M{"$set": M{
		"status":  models.FocusStatusHandled,
		"handled": models.Now().Unix(),
	}}
	return models.EmptyFocus.FindAndModify(query, change)
}

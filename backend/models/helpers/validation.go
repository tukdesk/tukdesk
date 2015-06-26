package helpers

import (
	"fmt"

	"github.com/tukdesk/tukdesk/backend/models"

	"github.com/tukdesk/httputils/validation"
	"gopkg.in/mgo.v2/bson"
)

const (
	msgForFocusResourceId = "is expected to be an ObjectId"
)

func ValidationNew() *validation.Validation {
	return validation.New()
}

// brand
func ValidationForBrandName(v *validation.Validation, key, name string) *validation.Validation {
	v.Required(key, name)
	v.MaxSize(key, name, BrandNameMaxLength)
	return v
}

// email
func ValidationForEmail(v *validation.Validation, key, email string) *validation.Validation {
	v.Required(key, email)
	v.Email(key, email)
	v.MaxSize(key, email, LimitedDataFieldMaxLength)
	return v
}

// user
func ValidationForUserPassword(v *validation.Validation, key, password string) *validation.Validation {
	v.Required(key, password)
	v.MinSize(key, password, UserPasswordMinLength)
	v.MaxSize(key, password, LimitedDataFieldMaxLength)
	return v
}

func ValidationForUserName(v *validation.Validation, key, name string) *validation.Validation {
	v.Required(key, name)
	v.MaxSize(key, name, UserNameMaxLength)
	return v
}

func ValidationForUserImportance(v *validation.Validation, key string, importance models.TypeUserImportance) *validation.Validation {
	v.In(key, importance, UserImportanceOptions)
	return v
}

func ValidationForUserListSort(v *validation.Validation, key, sort string) *validation.Validation {
	v.In(key, sort, UserSortOptionsForList)
	return v
}

// ticket
func ValidationForTicketSubject(v *validation.Validation, key, subject string) *validation.Validation {
	v.Required(key, subject)
	v.MaxSize(key, subject, TicketSubjectMaxLength)
	return v
}

func ValidationForTicketContent(v *validation.Validation, key, content string) *validation.Validation {
	v.Required(key, content)
	v.MinSize(key, content, TicketContentMinLength)
	return v
}

func ValidationForTicketChannel(v *validation.Validation, key, channel string) *validation.Validation {
	v.Required(key, channel)
	v.In(key, channel, TicketChannelOptionsForCreate)
	return v
}

func ValidationForTicketExtendField(v *validation.Validation, extend map[string]string) *validation.Validation {
	for key, val := range extend {
		v.MaxSize(fmt.Sprintf("extend.%s", key), val, TicketExtendFieldMaxLength)
	}
	return v
}

func ValidationForTicketPriority(v *validation.Validation, key string, priority models.TypePriority) *validation.Validation {
	v.In(key, priority, TicketPriorityOptions)
	return v
}

func ValidationForTicketStatusOnCreate(v *validation.Validation, key, status string) *validation.Validation {
	v.In(key, status, TicketStatusOptionsForCreate)
	return v
}

func ValidationForTicketStatusOnUpdate(v *validation.Validation, key, status string) *validation.Validation {
	v.In(key, status, TicketStatusOptionsForUpdate)
	return v
}

func ValidationForTicektListSort(v *validation.Validation, key, sort string) *validation.Validation {
	v.In(key, sort, TicketSortOptionsForList)
	return v
}

// comment
func ValidationForCommentTypeOnCreate(v *validation.Validation, key, typeName string) *validation.Validation {
	v.Required(key, typeName)
	v.In(key, typeName, CommentTypeOptionsForCreate)
	return v
}

func ValidationForCommentContent(v *validation.Validation, key, content string) *validation.Validation {
	v.Required(key, content)
	v.MinSize(key, content, CommentContentMinLength)
	return v
}

func ValidationForCommentTypeOnUpdate(v *validation.Validation, key, typeName string) *validation.Validation {
	v.Required(key, typeName)
	v.In(key, typeName, CommentTypeOptionsForUpdate)
	return v
}

// resource
func ValidationForFocusResourceId(v *validation.Validation, key string, resourceId bson.ObjectId) *validation.Validation {
	if IsEmptyId(resourceId) {
		v.AddError(key, msgForFocusResourceId)
	}
	return v
}

func ValidationForFocusResourceType(v *validation.Validation, key, resourceType string) *validation.Validation {
	v.In(key, resourceType, ResourceTypeOptions)
	return v
}

// focus
func ValidationForFoucsMessage(v *validation.Validation, key, message string) *validation.Validation {
	v.MaxSize(key, message, FocusMessageMaxLength)
	return v
}

func ValidationForFocusPriority(v *validation.Validation, key string, priority models.TypePriority) *validation.Validation {
	v.In(key, priority, FocusPriorityOptions)
	return v
}

func ValidationForFocusStatusOnList(v *validation.Validation, key, status string) *validation.Validation {
	v.In(key, status, FocusStatusOptionForList)
	return v
}

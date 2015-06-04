package helpers

import (
	"fmt"

	"github.com/tukdesk/httputils/validation"
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

func ValidationForBrandLogo(v *validation.Validation, key, logo string) *validation.Validation {
	v.MaxSize(key, logo, LimitedDataFieldMaxLength)
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

func ValidationForUserAvatar(v *validation.Validation, key, avatar string) *validation.Validation {
	v.MaxSize(key, avatar, LimitedDataFieldMaxLength)
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

func ValidationForTicketExtendField(v *validation.Validation, extend map[string]string) *validation.Validation {
	for key, val := range extend {
		v.MaxSize(fmt.Sprintf("extend.%s", key), val, TicketExtendFieldMaxLength)
	}
	return v
}

func ValidationForTicketStatusOnCreate(v *validation.Validation, key, status string) *validation.Validation {
	v.In(key, status, TicketStatusOptionsForCreate)
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

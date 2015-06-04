package apis

import (
	"fmt"
	"net/http"

	"github.com/tukdesk/httputils/jsonutils"
)

const (
	StatusCodeInternalError = 579

	ErrCodeInvalidRequestBody    = 990101
	ErrCodeInvalidArgs           = 990102
	ErrCodeInvalidJsonObjectType = 990105
)

func ErrorInvalidRequestBodyWithError(err error) error {
	return jsonutils.NewAPIError(http.StatusBadRequest, ErrCodeInvalidRequestBody, err.Error())
}

func ErrorInvaidArgsWithMsg(msg string) error {
	return jsonutils.NewAPIError(http.StatusBadRequest, ErrCodeInvalidArgs, msg)
}

func ErrorInvalidArgType(name, typeName string) error {
	return jsonutils.NewAPIError(http.StatusBadRequest, ErrCodeInvalidJsonObjectType, fmt.Sprintf("%s is expected to be %s", name, typeName))
}

var (
	ErrInternalError = jsonutils.NewAPIError(StatusCodeInternalError, 990103, "internal error")
	ErrInvalidId     = jsonutils.NewAPIError(http.StatusBadRequest, 990104, "invalid id")

	ErrUnauthorized = jsonutils.NewAPIError(http.StatusForbidden, 990201, "unauthorized request")
	ErrUnlogged     = jsonutils.NewAPIError(http.StatusForbidden, 990202, "unlogged")
	ErrAgentOnly    = jsonutils.NewAPIError(http.StatusForbidden, 990203, "for agent only")

	ErrBrandAlreadyInitialized = jsonutils.NewAPIError(http.StatusForbidden, 110101, "brand already initialized")
	ErrBrandNotFound           = jsonutils.NewAPIError(http.StatusBadRequest, 110102, "brand not found")

	ErrAgentNotFound         = jsonutils.NewAPIError(http.StatusBadRequest, 110201, "agent not found")
	ErrAgentPasswordNotMatch = jsonutils.NewAPIError(http.StatusForbidden, 110202, "passwod not match")

	ErrTicketDuplicate = jsonutils.NewAPIError(http.StatusBadRequest, 110301, "ticket with duplicate channel")
	ErrTicketNotFound  = jsonutils.NewAPIError(http.StatusNotFound, 110302, "ticket not found")

	ErrCommentNotFound     = jsonutils.NewAPIError(http.StatusNotFound, 110401, "comment not found")
	ErrCommentUnchangeable = jsonutils.NewAPIError(http.StatusForbidden, 110402, "comment unchangeable")

	ErrUserNotFound = jsonutils.NewAPIError(http.StatusNotFound, 110501, "user not found")

	ErrFocusNotFound       = jsonutils.NewAPIError(http.StatusNotFound, 110601, "focus not found")
	ErrFocusAlreadyHandled = jsonutils.NewAPIError(http.StatusBadRequest, 110602, "focus already handled")

	ErrResourceNotFound = jsonutils.NewAPIError(http.StatusNotFound, 110701, "resource not found")
)

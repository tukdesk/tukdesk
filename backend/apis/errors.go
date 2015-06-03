package apis

import (
	"net/http"

	"github.com/tukdesk/httputils/jsonutils"
)

const (
	StatusCodeInternalError = 579

	ErrCodeInvalidRequestBody = 990101
	ErrCodeInvalidArgs        = 990102
)

func ErrInvalidRequestBodyWithError(err error) error {
	return jsonutils.NewAPIError(http.StatusBadRequest, ErrCodeInvalidRequestBody, err.Error())
}

func ErrInvaidArgsWithMsg(msg string) error {
	return jsonutils.NewAPIError(http.StatusBadRequest, ErrCodeInvalidArgs, msg)
}

var (
	ErrInternalError = jsonutils.NewAPIError(StatusCodeInternalError, 990103, "internal error")
	ErrUnauthorized  = jsonutils.NewAPIError(http.StatusForbidden, 990201, "unauthorized request")
	ErrUnlogged      = jsonutils.NewAPIError(http.StatusForbidden, 990202, "unlogged")

	ErrBrandAlreadyInitialized = jsonutils.NewAPIError(http.StatusForbidden, 110101, "brand already initialized")
	ErrBrandNotFound           = jsonutils.NewAPIError(http.StatusBadRequest, 110102, "brand not found")

	ErrAgentNotFound         = jsonutils.NewAPIError(http.StatusBadRequest, 110201, "agent not found")
	ErrAgentPasswordNotMatch = jsonutils.NewAPIError(http.StatusForbidden, 110202, "passwod not match")

	ErrTicketDuplicate = jsonutils.NewAPIError(http.StatusBadRequest, 110301, "ticket with duplicate channel")
)

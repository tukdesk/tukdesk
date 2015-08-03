package apis

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
	emw "github.com/tukdesk/httputils/echomiddleware"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/tukdesk/httputils/validation"
	"github.com/tukdesk/httputils/xlogger"

	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

const (
	currentUserKey = "_user"
)

func abort(err error) {
	panic(err)
	return
}

func NowUnix() int64 {
	return time.Now().Unix()
}

func GetCurrentBrand() *models.Brand {
	brand := helpers.CurrentBrand()
	if brand == nil {
		abort(ErrBrandNotFound)
		return nil
	}

	return brand
}

func GetCurrentUser(c *echo.Context) *models.User {
	user, _ := c.Get(currentUserKey).(*models.User)
	return user
}

func CheckAuthorizedAsAgent(c *echo.Context) *models.User {
	user := GetCurrentUser(c)
	if !AuthorizedAsAgent(user) {
		abort(ErrAgentOnly)
	}
	return user
}

func CheckAuthorizedLogged(c *echo.Context) *models.User {
	user := GetCurrentUser(c)
	if !AuthorizedLogged(user) {
		abort(ErrUnlogged)
	}
	return user
}

func CheckValidation(v *validation.Validation) {
	if !v.HasErrors() {
		return
	}

	e := v.Errors()[0]
	abort(ErrorInvaidArgsWithMsg(fmt.Sprintf("%s : %s", e.Key, e.Message)))
	return
}

func GetJsonArgsFromContext(c *echo.Context, args interface{}) {
	if err := jsonutils.GetJsonArgsFromRequest(c.Request(), args); err != nil {
		abort(ErrorInvalidRequestBodyWithError(err))
		return
	}
}

func GetMapArgsFromContext(c *echo.Context) map[string]interface{} {
	m := map[string]interface{}{}
	if err := jsonutils.GetJsonArgsFromRequest(c.Request(), &m); err != nil {
		abort(ErrorInvalidRequestBodyWithError(err))
		return nil
	}
	return m
}

func GetLogger(c *echo.Context) *xlogger.XLogger {
	return emw.GetRequestLogger(c)
}

func OutputJson(data interface{}, c *echo.Context) error {
	return c.JSON(StatusCodeOK, data)
}

func ChangeSetM(setM map[string]interface{}) helpers.M {
	return helpers.M{"$set": setM}
}

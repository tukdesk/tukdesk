package apis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
	"github.com/zenazn/goji/web"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/tukdesk/httputils/validation"
	"github.com/tukdesk/httputils/xlogger"
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

func CheckAuthorizedAsAgent(c *web.C, w http.ResponseWriter, r *http.Request) *models.User {
	user := GetCurrentUser(c, w, r)
	if !AuthorizedAsAgent(user) {
		abort(ErrAgentOnly)
	}
	return user
}

func CheckAuthorizedLogged(c *web.C, w http.ResponseWriter, r *http.Request) *models.User {
	user := GetCurrentUser(c, w, r)
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

func GetJsonArgsFromRequest(r *http.Request, args interface{}) {
	if err := jsonutils.GetJsonArgsFromRequest(r, args); err != nil {
		abort(ErrorInvalidRequestBodyWithError(err))
		return
	}
}

func GetMapArgsFromRequest(r *http.Request) map[string]interface{} {
	m := map[string]interface{}{}
	if err := jsonutils.GetJsonArgsFromRequest(r, &m); err != nil {
		abort(ErrorInvalidRequestBodyWithError(err))
		return nil
	}
	return m
}

func GetLogger(c *web.C, w http.ResponseWriter, r *http.Request) *xlogger.XLogger {
	return gojimiddleware.GetRequestLogger(c, w, r)
}

func OutputJson(data interface{}, w http.ResponseWriter, r *http.Request) {
	jsonutils.OutputJson(data, w, r)
	return
}

func ChangeSetM(setM map[string]interface{}) helpers.M {
	return helpers.M{"$set": setM}
}

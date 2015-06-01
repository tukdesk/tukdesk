package apis

import (
	"fmt"
	"net/http"

	"github.com/tukdesk/tukdesk/backend/models/helpers"
	"github.com/zenazn/goji/web"

	"github.com/astaxie/beego/validation"
	"github.com/tukdesk/httputils/jsonutils"
)

func abort(err error) {
	panic(err)
	return
}

func CheckCurrentBrand() {
	if helpers.CurrentBrand() == nil {
		abort(ErrBrandNotFound)
	}

	return
}

func CheckAuthorizedAsAgent(c *web.C) {
	user := GetCurrentUser(c)
	if !AuthorizedAsAgent(user) {
		abort(ErrUnauthorized)
	}
	return
}

func FirstError(v *validation.Validation) error {
	if !v.HasErrors() {
		return nil
	}

	e := v.Errors[0]
	return ErrInvaidArgsWithMsg(fmt.Sprintf("%s : %s", e.Key, e.Message))
}

func GetJsonArgsFromRequest(r *http.Request, args interface{}) {
	if err := jsonutils.GetJsonArgsFromRequest(r, args); err != nil {
		abort(ErrInvalidRequestBodyWithError(err))
		return
	}
}

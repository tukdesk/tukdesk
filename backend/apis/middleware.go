package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/zenazn/goji/web"
)

const (
	currentUserKey = "_user"
)

func CurrentUser(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if helpers.CurrentBrand() == nil {
			jsonutils.OutputJsonError(ErrBrandNotFound, w, r)
			return
		}

		user, _, err := helpers.UserFromRequest(r, helpers.CurrentBrand().APIKey)
		if err != nil && err != helpers.ErrTokenNotFound && !helpers.IsInvalidToken(err) {
			logger := gojimiddleware.GetRequestLogger(c, w, r)
			logger.Error(err)
			jsonutils.OutputJsonError(ErrInternalError, w, r)
			return
		}

		c.Env[currentUserKey] = user
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetCurrentUser(c *web.C) *models.User {
	return c.Env[currentUserKey].(*models.User)
}

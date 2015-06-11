package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

const (
	currentUserKey = "_user"
)

func CurrentUser(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// get current user
		GetCurrentUser(c, w, r)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetCurrentUser(c *web.C, w http.ResponseWriter, r *http.Request) *models.User {
	CheckCurrentBrand()

	if user, ok := c.Env[currentUserKey].(*models.User); ok {
		return user
	}

	user, _, err := helpers.UserFromRequest(r, helpers.CurrentBrand().Authorization.APIKey)
	if err != nil && err != helpers.ErrTokenNotFound && !helpers.IsInvalidToken(err) {
		logger := gojimiddleware.GetRequestLogger(c, w, r)
		logger.Error(err)
		abort(ErrInternalError)
		return nil
	}

	if err != nil {
		user = nil
	}

	c.Env[currentUserKey] = user
	return user
}

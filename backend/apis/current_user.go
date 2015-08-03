package apis

import (
	"github.com/labstack/echo"

	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

func CurrentUser(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		user, err := GetCurrentUserFromContext(c)

		if err != nil {
			logger := GetLogger(c)
			logger.Error(err)
			return ErrInternalError
		}

		c.Set(currentUserKey, user)

		defer func() {
			c.Set(currentUserKey, nil)
		}()

		return h(c)
	}
}

func GetCurrentUserFromContext(c *echo.Context) (*models.User, error) {
	brand := GetCurrentBrand()

	user, _, err := helpers.UserFromRequest(c.Request(), brand.Authorization.APIKey)
	if err != nil && err != helpers.ErrTokenNotFound && !helpers.IsInvalidToken(err) {
		return nil, err
	}

	// token not found or invalid token
	if err != nil {
		return nil, nil
	}

	return user, nil
}

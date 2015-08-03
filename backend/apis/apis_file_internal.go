package apis

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tukdesk/httputils/tools"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

const (
	defaultMultipartMaxMemory = 4 << 20 // 4M
)

type InternalFileModule struct {
	cfg      *config.Config
	storager FileStorager
}

func newInternalFileModule(cfg *config.Config, group *echo.Group) error {
	if cfg.File.Internal.Dir == "" {
		return fmt.Errorf("file dir required")
	}

	storager, err := newInternalLocalStorager(cfg.File.Internal.Dir)
	if err != nil {
		return err
	}

	m := InternalFileModule{
		cfg:      cfg,
		storager: storager,
	}

	group.Get("/token", m.token)
	group.Post("/upload", m.upload)
	return nil
}

func (this *InternalFileModule) upload(c *echo.Context) error {
	brand := GetCurrentBrand()
	user := GetCurrentUser(c)
	// check file token
	logger := GetLogger(c)

	r := c.Request()

	if err := r.ParseMultipartForm(defaultMultipartMaxMemory); err != nil {
		logger.Error(err)
		return ErrFileInternalInvalidRequest
	}

	if !helpers.InternalFileTokenValid(getUserIdentifier(r, user), multipartFormValue(r, "token"), brand.Authorization.APIKey) {
		return ErrFileInternalInvalidToken
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return ErrFileInternalFileNotFound
	}

	fileHeaders := r.MultipartForm.File["file"]
	if len(fileHeaders) == 0 {
		return ErrFileInternalFileNotFound
	}

	fileHeader := fileHeaders[0]
	fileDoc, err := this.storager.Store(fileHeader)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	if err := fileDoc.Insert(); err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(fileDoc, c)
}

func (this *InternalFileModule) token(c *echo.Context) error {
	brand := GetCurrentBrand()
	user := GetCurrentUser(c)

	userIdentifier := getUserIdentifier(c.Request(), user)

	output := helpers.OutputTokenInfo(helpers.NewInternalFileToken(userIdentifier, brand.Authorization.APIKey, helpers.FileTokenExpiration), helpers.FileTokenExpirationSec)
	return OutputJson(output, c)
}

func getUserIdentifier(r *http.Request, user *models.User) string {
	if AuthorizedLogged(user) {
		return user.Id.Hex()
	}
	return tools.GetRealIp(r)
}

func multipartFormValue(r *http.Request, key string) string {
	if r.MultipartForm == nil || r.MultipartForm.Value == nil {
		return ""
	}

	vals := r.MultipartForm.Value[key]
	if len(vals) == 0 {
		return ""
	}

	return vals[0]
}

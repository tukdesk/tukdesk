package apis

import (
	"fmt"
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/tools"
	"github.com/zenazn/goji/web"
)

const (
	defaultMultipartMaxMemory = 4 << 20 // 4M
)

type InternalFileModule struct {
	cfg      *config.Config
	storager FileStorager
}

func newInternalFileModule(cfg *config.Config) (*web.Mux, error) {
	if cfg.File.Internal.Dir == "" {
		return nil, fmt.Errorf("file dir required")
	}

	storager, err := newInternalLocalStorager(cfg.File.Internal.Dir)
	if err != nil {
		return nil, err
	}

	m := InternalFileModule{
		cfg:      cfg,
		storager: storager,
	}

	mux := web.New()
	mux.Get("/token", m.token)
	mux.Post("/upload", m.upload)
	return mux, nil
}

func (this *InternalFileModule) upload(c web.C, w http.ResponseWriter, r *http.Request) {
	brand := GetCurrentBrand()
	user := GetCurrentUser(&c, w, r)
	// check file token
	logger := GetLogger(&c, w, r)

	if err := r.ParseMultipartForm(defaultMultipartMaxMemory); err != nil {
		logger.Error(err)
		abort(ErrFileInternalInvalidRequest)
		return
	}

	if !helpers.InternalFileTokenValid(getUserIdentifier(r, user), multipartFormValue(r, "token"), brand.Authorization.APIKey) {
		abort(ErrFileInternalInvalidToken)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		abort(ErrFileInternalFileNotFound)
		return
	}

	fileHeaders := r.MultipartForm.File["file"]
	if len(fileHeaders) == 0 {
		abort(ErrFileInternalFileNotFound)
		return
	}

	fileHeader := fileHeaders[0]
	fileDoc, err := this.storager.Store(fileHeader)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	if err := fileDoc.Insert(); err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(fileDoc, w, r)
	return
}

func (this *InternalFileModule) token(c web.C, w http.ResponseWriter, r *http.Request) {
	brand := GetCurrentBrand()
	user := GetCurrentUser(&c, w, r)

	userIdentifier := getUserIdentifier(r, user)

	output := helpers.OutputTokenInfo(helpers.NewInternalFileToken(userIdentifier, brand.Authorization.APIKey, helpers.FileTokenExpiration), helpers.FileTokenExpirationSec)
	OutputJson(output, w, r)
	return
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

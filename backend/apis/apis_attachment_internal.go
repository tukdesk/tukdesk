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

type InternalAttachmentModule struct {
	cfg      *config.Config
	storager AttachmentStorager
}

func newInternalAttachmentModule(cfg *config.Config) (*web.Mux, error) {
	if cfg.Attachment.Internal.Dir == "" {
		return nil, fmt.Errorf("attachment dir required")
	}

	storager, err := newInternalLocalStorager(cfg.Attachment.Internal.Dir)
	if err != nil {
		return nil, err
	}

	m := InternalAttachmentModule{
		cfg:      cfg,
		storager: storager,
	}

	mux := web.New()
	mux.Get("/token", m.token)
	mux.Post("/upload", m.upload)
	return mux, nil
}

func (this *InternalAttachmentModule) upload(c web.C, w http.ResponseWriter, r *http.Request) {
	brand := GetCurrentBrand()
	user := GetCurrentUser(&c, w, r)
	// check attachment token
	logger := GetLogger(&c, w, r)

	if err := r.ParseMultipartForm(defaultMultipartMaxMemory); err != nil {
		logger.Error(err)
		abort(ErrAttachmentInternalInvalidRequest)
		return
	}

	if !helpers.InternalAttachmentTokenValid(getUserIdentifier(r, user), multipartFormValue(r, "token"), brand.Authorization.APIKey) {
		abort(ErrAttachmentInternalInvalidToken)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		abort(ErrAttachmentInternalFileNotFound)
		return
	}

	fileHeaders := r.MultipartForm.File["file"]
	if len(fileHeaders) == 0 {
		abort(ErrAttachmentInternalFileNotFound)
		return
	}

	fileHeader := fileHeaders[0]
	attachment, err := this.storager.Store(fileHeader)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	if err := attachment.Insert(); err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(attachment, w, r)
	return
}

func (this *InternalAttachmentModule) token(c web.C, w http.ResponseWriter, r *http.Request) {
	brand := GetCurrentBrand()
	user := GetCurrentUser(&c, w, r)

	userIdentifier := getUserIdentifier(r, user)

	output := helpers.OutputTokenInfo(helpers.NewInternalAttachmentToken(userIdentifier, brand.Authorization.APIKey, helpers.AttachmentTokenExpiration), helpers.AttachmentTokenExpirationSec)
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

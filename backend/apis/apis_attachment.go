package apis

import (
	"github.com/tukdesk/tukdesk/backend/config"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

func RegisterAttachmentModule(cfg *config.Config, app *web.Mux) (*web.Mux, error) {
	mux, err := newInternalAttachmentModule(cfg)
	if err != nil {
		return nil, err
	}

	gojimiddleware.RegisterSubroute("/attachments", app, mux)
	return mux, nil
}

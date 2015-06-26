package apis

import (
	"github.com/tukdesk/tukdesk/backend/config"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

func RegisterFileModule(cfg *config.Config, app *web.Mux) (*web.Mux, error) {
	mux, err := newInternalFileModule(cfg)
	if err != nil {
		return nil, err
	}

	gojimiddleware.RegisterSubroute("/files", app, mux)
	return mux, nil
}

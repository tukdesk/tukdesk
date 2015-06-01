package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/zenazn/goji/web"
)

type BrandModule struct {
	cfg config.Config
}

func RegisterBrandModule(cfg config.Config, app *web.Mux) *web.Mux {
	m := BrandModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Use(CurrentUser)
	mux.Get("/key", m.brandAPIKey)

	gojimiddleware.RegisterSubroute("/brand", app, mux)
	return mux
}

func (this *BrandModule) brandAPIKey(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c)
	logger := gojimiddleware.GetRequestLogger(&c, w, r)
	logger.Info(GetCurrentUser(&c))
	output := &helpers.OutputAPIKey{
		Key: string(helpers.CurrentBrand().APIKey),
	}
	jsonutils.OutputJson(output, w, r)
	return
}

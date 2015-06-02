package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/astaxie/beego/validation"
	"github.com/tukdesk/httputils/gojimiddleware"
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
	mux.Get("", m.brandInfo)
	mux.Put("", m.brandUpdate)
	mux.Get("/key", m.brandAPIKey)
	mux.Put("/key", m.brandResetAPIKey)

	gojimiddleware.RegisterSubroute("/brand", app, mux)
	return mux
}

func (this *BrandModule) brandAPIKey(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c, w, r)
	output := helpers.OutputAPIKey(helpers.CurrentBrand().Authorization.APIKey)
	OutputJson(output, w, r)
	return
}

func (this *BrandModule) brandInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	OutputJson(helpers.CurrentBrand(), w, r)
	return
}

func (this *BrandModule) brandUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c, w, r)

	args := &BrandUpdateArgs{}
	GetJsonArgsFromRequest(r, args)

	v := &validation.Validation{}
	setM := helpers.M{}
	brand := helpers.CurrentBrand()

	if args.Base.Name != "" && args.Base.Name != brand.Base.Name {
		v.MaxSize(args.Base.Name, helpers.BrandNameMaxLength, "name")
		setM["base.name"] = args.Base.Name
	}

	if args.Base.Logo != "" && args.Base.Logo != brand.Base.Logo {
		v.MaxSize(args.Base.Logo, helpers.LimitedDataFieldMaxLength, "logo")
		setM["base.logo"] = args.Base.Logo
	}

	CheckValidation(v)

	if len(setM) > 0 {
		setM["updated"] = NowUnix()
		if err := helpers.BrandUpdateCurrent(ChangeSetM(setM)); err != nil {
			GetLogger(&c, w, r).Error(err)
			abort(ErrInternalError)
			return
		}

	}

	OutputJson(brand, w, r)
	return
}

func (this *BrandModule) brandResetAPIKey(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c, w, r)

	newApiKey := helpers.BrandNewAPIKey()
	setM := helpers.M{
		"authorization.apiKey": newApiKey,
		"updated":              NowUnix(),
	}

	if err := helpers.BrandUpdateCurrent(ChangeSetM(setM)); err != nil {
		GetLogger(&c, w, r).Error(err)
		abort(ErrInternalError)
		return
	}

	output := helpers.OutputAPIKey(helpers.CurrentBrand().Authorization.APIKey)
	OutputJson(output, w, r)
	return
}

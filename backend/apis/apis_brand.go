package apis

import (
	"github.com/labstack/echo"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

type BrandModule struct {
	cfg *config.Config
}

func RegisterBrandModule(cfg *config.Config, mux *echo.Group) {
	m := BrandModule{
		cfg: cfg,
	}

	group := mux.Group("/brand")
	group.Use(CurrentUser)

	group.Get("", m.brandInfo)
	group.Put("", m.brandUpdate)
	group.Get("/key", m.brandAPIKey)
	group.Put("/key", m.brandResetAPIKey)
	return
}

func (this *BrandModule) brandAPIKey(c *echo.Context) error {
	CheckAuthorizedAsAgent(c)
	output := helpers.OutputBrandAPIKey(helpers.CurrentBrand().Authorization.APIKey)
	return OutputJson(output, c)
}

func (this *BrandModule) brandInfo(c *echo.Context) error {
	return OutputJson(helpers.CurrentBrand(), c)
}

func (this *BrandModule) brandUpdate(c *echo.Context) error {
	CheckAuthorizedAsAgent(c)

	args := &BrandUpdateArgs{}
	GetJsonArgsFromContext(c, args)

	v := helpers.ValidationNew()
	setM := helpers.M{}
	brand := helpers.CurrentBrand()

	if args.Base.Name != "" && args.Base.Name != brand.Base.Name {
		helpers.ValidationForBrandName(v, "name", args.Base.Name)
		setM["base.name"] = args.Base.Name
	}

	if args.Base.Logo != nil {
		setM["base.logo"] = args.Base.Logo
	}

	CheckValidation(v)

	if len(setM) > 0 {
		setM["updated"] = NowUnix()
		if err := helpers.BrandUpdateCurrent(ChangeSetM(setM)); err != nil {
			GetLogger(c).Error(err)
			return ErrInternalError
		}

	}

	return OutputJson(brand, c)
}

func (this *BrandModule) brandResetAPIKey(c *echo.Context) error {
	CheckAuthorizedAsAgent(c)

	newApiKey := helpers.BrandNewAPIKey()
	setM := helpers.M{
		"authorization.apiKey": newApiKey,
		"updated":              NowUnix(),
	}

	if err := helpers.BrandUpdateCurrent(ChangeSetM(setM)); err != nil {
		GetLogger(c).Error(err)
		return ErrInternalError
	}

	output := helpers.OutputBrandAPIKey(helpers.CurrentBrand().Authorization.APIKey)
	return OutputJson(output, c)
}

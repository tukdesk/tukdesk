package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

type BaseModule struct {
	cfg *config.Config
}

func RegisterBaseModule(cfg *config.Config, app *web.Mux) *web.Mux {
	m := BaseModule{
		cfg: cfg,
	}

	mux := web.New()

	mux.Post("/init", m.brandInit)
	mux.Post("/signin", m.signin)

	gojimiddleware.RegisterSubroute("/base", app, mux)
	return mux
}

func (this *BaseModule) brandInit(c web.C, w http.ResponseWriter, r *http.Request) {
	if helpers.CurrentBrand() != nil {
		abort(ErrBrandAlreadyInitialized)
		return
	}

	args := &BrandInitArgs{}
	GetJsonArgsFromRequest(r, args)

	v := helpers.ValidationNew()
	helpers.ValidationForBrandName(v, "brand name", args.BrandName)
	if args.Name != "" {
		helpers.ValidationForUserName(v, "name", args.Name)
	}
	helpers.ValidationForEmail(v, "email", args.Email)
	helpers.ValidationForUserPassword(v, "password", args.Password)

	CheckValidation(v)

	logger := GetLogger(&c, w, r)

	// brand init
	brand, err := helpers.BrandInit(args.BrandName)
	if helpers.IsDup(err) {
		abort(ErrBrandAlreadyInitialized)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	// user init
	if args.Name == "" {
		args.Name = helpers.UserGetValidNameFromEmail(args.Email)
	}

	_, err = helpers.AgentInit(args.Email, args.Name, args.Password, this.cfg.Salt)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(brand, w, r)
	return
}

func (this *BaseModule) signin(c web.C, w http.ResponseWriter, r *http.Request) {
	brand := GetCurrentBrand()

	args := &SigninArgs{}
	GetJsonArgsFromRequest(r, args)

	v := helpers.ValidationNew()
	helpers.ValidationForUserPassword(v, "password", args.Password)

	CheckValidation(v)

	logger := GetLogger(&c, w, r)

	user, err := helpers.AgentFind()
	if helpers.IsNotFound(err) {
		abort(ErrAgentNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	if !helpers.UserCheckPassword(user, args.Password, this.cfg.Salt) {
		abort(ErrAgentPasswordNotMatch)
		return
	}

	output := helpers.OutputTokenInfo(helpers.TokenForUser(user, brand.Authorization.APIKey), helpers.TokenDefaultExpirationSec)
	OutputJson(output, w, r)
	return
}

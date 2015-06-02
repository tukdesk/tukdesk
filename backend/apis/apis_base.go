package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/astaxie/beego/validation"
	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

type BaseModule struct {
	cfg config.Config
}

func RegisterBaseModule(cfg config.Config, app *web.Mux) *web.Mux {
	m := BaseModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Post("/init", m.brandInit)
	mux.Post("/signup", m.signup)

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

	v := &validation.Validation{}
	v.Required(args.BrandName, "brand name")
	v.MaxSize(args.BrandName, helpers.BrandNameMaxLength, "brand name")
	v.MaxSize(args.Name, helpers.UserNameMaxLength, "name")
	v.Required(args.Email, "email")
	v.Email(args.Email, args.Email)
	v.Required(args.Password, "password")
	v.MinSize(args.Password, helpers.UserPasswordMinLength, "password")

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

	_, err = helpers.AgentInit(args.Email, args.Name, args.Password, brand.Authorization.Salt)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(brand, w, r)
	return
}

func (this *BaseModule) signup(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckCurrentBrand()

	args := &SignupArgs{}
	GetJsonArgsFromRequest(r, args)

	v := &validation.Validation{}

	v.Required(args.Password, "password")
	v.MinSize(args.Password, helpers.UserPasswordMinLength, "password")
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

	if !helpers.UserCheckPassword(user, args.Password, helpers.CurrentBrand().Authorization.Salt) {
		abort(ErrAgentPasswordNotMatch)
		return
	}

	output := helpers.OuputToken(helpers.TokenForUser(user, helpers.CurrentBrand().Authorization.APIKey), helpers.TokenDefaultExpirationSec)
	OutputJson(output, w, r)
	return
}

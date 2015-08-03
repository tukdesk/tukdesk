package apis

import (
	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/labstack/echo"
)

type BaseModule struct {
	cfg *config.Config
}

func RegisterBaseModule(cfg *config.Config, mux *echo.Group) {
	m := BaseModule{
		cfg: cfg,
	}

	group := mux.Group("/base")
	group.Post("/init", m.brandInit)
	group.Post("/signin", m.signin)

	return
}

func (this *BaseModule) brandInit(c *echo.Context) error {
	if helpers.CurrentBrand() != nil {
		return ErrBrandAlreadyInitialized
	}

	args := &BrandInitArgs{}
	GetJsonArgsFromContext(c, args)

	v := helpers.ValidationNew()
	helpers.ValidationForBrandName(v, "brand name", args.BrandName)
	if args.Name != "" {
		helpers.ValidationForUserName(v, "name", args.Name)
	}
	helpers.ValidationForEmail(v, "email", args.Email)
	helpers.ValidationForUserPassword(v, "password", args.Password)

	CheckValidation(v)

	logger := GetLogger(c)

	// brand init
	brand, err := helpers.BrandInit(args.BrandName)
	if helpers.IsDup(err) {
		return ErrBrandAlreadyInitialized
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	// user init
	if args.Name == "" {
		args.Name = helpers.UserGetValidNameFromEmail(args.Email)
	}

	_, err = helpers.AgentInit(args.Email, args.Name, args.Password, this.cfg.Salt)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return c.JSON(StatusCodeOK, brand)
}

func (this *BaseModule) signin(c *echo.Context) error {
	brand := GetCurrentBrand()

	args := &SigninArgs{}
	GetJsonArgsFromContext(c, args)

	v := helpers.ValidationNew()
	helpers.ValidationForUserPassword(v, "password", args.Password)

	CheckValidation(v)

	logger := GetLogger(c)

	user, err := helpers.AgentFind()
	if helpers.IsNotFound(err) {
		return ErrAgentNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	if !helpers.UserCheckPassword(user, args.Password, this.cfg.Salt) {
		return ErrAgentPasswordNotMatch
	}

	output := helpers.OutputTokenInfo(helpers.TokenForUser(user, brand.Authorization.APIKey), helpers.TokenDefaultExpirationSec)
	return c.JSON(StatusCodeOK, output)
}

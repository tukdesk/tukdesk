package apis

import (
	"github.com/labstack/echo"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

type ProfileModule struct {
	cfg *config.Config
}

func RegisterProfileModule(cfg *config.Config, mux *echo.Group) {
	m := ProfileModule{
		cfg: cfg,
	}

	group := mux.Group("/profile")
	group.Use(CurrentUser)

	group.Get("", m.profile)
	group.Put("", m.profileUpdate)
	return
}

func (this *ProfileModule) profile(c *echo.Context) error {
	user := CheckAuthorizedLogged(c)
	output := helpers.OutputUserProfileInfo(user)
	return OutputJson(output, c)
}

func (this *ProfileModule) profileUpdate(c *echo.Context) error {
	user := CheckAuthorizedLogged(c)

	args := &ProfileUpdateArgs{}
	GetJsonArgsFromContext(c, args)

	v := helpers.ValidationNew()
	setM := helpers.M{}

	if args.Base.Name != "" && args.Base.Name != user.Base.Name {
		helpers.ValidationForUserName(v, "name", args.Base.Name)
		setM["base.name"] = args.Base.Name
	}

	if args.Base.Avatar != nil {
		setM["base.avatar"] = args.Base.Avatar
	}

	CheckValidation(v)

	if len(setM) > 0 {
		setM["updated"] = NowUnix()
		if err := helpers.UserFindAndModify(user, ChangeSetM(setM)); err != nil {
			GetLogger(c).Error(err)
			return ErrInternalError
		}
	}

	output := helpers.OutputUserProfileInfo(user)
	return OutputJson(output, c)
}

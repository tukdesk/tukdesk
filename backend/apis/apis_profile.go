package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/astaxie/beego/validation"
	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

type ProfileModule struct {
	cfg config.Config
}

func RegisterProfileModule(cfg config.Config, app *web.Mux) *web.Mux {
	m := ProfileModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Get("", m.profile)
	mux.Put("", m.profileUpdate)

	gojimiddleware.RegisterSubroute("/profile", app, mux)
	return mux
}

func (this *ProfileModule) profile(c web.C, w http.ResponseWriter, r *http.Request) {
	user := CheckAuthorizedLogged(&c, w, r)
	output := helpers.OutputProfile(user)
	OutputJson(output, w, r)
	return
}

func (this *ProfileModule) profileUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	user := CheckAuthorizedLogged(&c, w, r)

	args := &ProfileUpdateArgs{}
	GetJsonArgsFromRequest(r, args)

	v := &validation.Validation{}
	setM := helpers.M{}

	if args.Base.Name != "" && args.Base.Name != user.Base.Name {
		v.MaxSize(args.Base.Name, helpers.UserNameMaxLength, "name")
		setM["base.name"] = args.Base.Name
	}

	if args.Base.Avatar != "" && args.Base.Avatar != user.Base.Avatar {
		v.MaxSize(args.Base.Avatar, helpers.LimitedDataFieldMaxLength, "avatar")
		setM["base.avatar"] = args.Base.Avatar
	}

	CheckValidation(v)

	if len(setM) > 0 {
		setM["updated"] = NowUnix()
		if err := helpers.UserFindAndModifyWithUser(user, ChangeSetM(setM)); err != nil {
			GetLogger(&c, w, r).Error(err)
			abort(ErrInternalError)
			return
		}
	}

	output := helpers.OutputProfile(user)
	OutputJson(output, w, r)
	return
}

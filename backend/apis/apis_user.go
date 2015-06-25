package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/zenazn/goji/web"
)

var (
	defaultClientsSort = []string{"-updated"}
)

type UserModule struct {
	cfg *config.Config
}

func RegisterUserModule(cfg *config.Config, app *web.Mux) *web.Mux {
	m := UserModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Get("", m.userList)
	mux.Get("/:userId", m.userInfo)
	mux.Put("/:userId", m.userUpdate)
	mux.Use(CurrentUser)

	gojimiddleware.RegisterSubroute("/users", app, mux)
	return mux
}

func (this *UserModule) userList(c web.C, w http.ResponseWriter, r *http.Request) {
	// 只有 agent 可以查看全员列表
	CheckAuthorizedAsAgent(&c, w, r)

	listArgs := GetListArgsFromRequest(r)
	v := helpers.ValidationNew()

	if listArgs.Sort != "" {
		helpers.ValidationForUserListSort(v, "sort", listArgs.Sort)
	}

	CheckValidation(v)

	query := helpers.M{}

	if importanceStr := r.FormValue("importance"); importanceStr != "" {
		importance := models.NewTypeUserImportance(importanceStr)
		helpers.ValidationForUserImportance(v, "importance", importance)
		CheckValidation(v)

		query["business.importance"] = importance
	}

	sort := make([]string, 0, len(defaultClientsSort)+1)
	if listArgs.Sort != "" {
		sort = append(sort, listArgs.Sort)
	}
	sort = append(sort, defaultClientsSort...)

	logger := GetLogger(&c, w, r)

	count, err := helpers.ClientCount(query)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	users, err := helpers.ClientListAfter(query, listArgs.LastId, listArgs.Limit, sort)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	items := make([]*helpers.OutputUser, len(users))
	if len(items) > 0 {
		for i, user := range users {
			items[i] = helpers.OutputUserProfileInfo(user)
		}
	}

	OutputJson(ListResultNew(count, items), w, r)
	return
}

func (this *UserModule) userInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(&c, w, r)

	userId, ok := helpers.IdFromString(c.URLParams["userId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	logger := GetLogger(&c, w, r)

	user, err := helpers.UserFindById(userId)
	if helpers.IsNotFound(err) {
		abort(ErrUserNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	var infoParser func(*models.User) *helpers.OutputUser
	if AuthorizedAsAgent(user) {
		infoParser = helpers.OutputUserDetailInfo
	} else {
		infoParser = helpers.OutputUserBaseInfo
	}

	OutputJson(infoParser(user), w, r)
	return
}

func (this *UserModule) userUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// 只有 agent 可以通过这个接口修改 user 信息
	// 只能修改 client 信息
	user := CheckAuthorizedAsAgent(&c, w, r)

	userId, ok := helpers.IdFromString(c.URLParams["userId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	args := GetMapArgsFromRequest(r)

	logger := GetLogger(&c, w, r)

	user, err := helpers.ClientFindById(userId)
	if helpers.IsNotFound(err) {
		abort(ErrUserNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	setM := helpers.M{}
	v := helpers.ValidationNew()

	for name, val := range args {
		switch name {
		case "business":
			businessMap, ok := val.(map[string]interface{})
			if !ok {
				abort(ErrorInvalidArgType(name, helpers.JSONTypeNameObject))
				return
			}

			if importanceStr, ok := businessMap["importance"].(string); ok {
				importance := models.NewTypeUserImportance(importanceStr)

				if importance == user.Business.Importance {
					break
				}

				helpers.ValidationForUserImportance(v, "business.importance", importance)
				CheckValidation(v)

				setM["business.importance"] = importance

			}
		}
	}

	if len(setM) > 0 {
		setM["updated"] = NowUnix()
		if err := helpers.UserFindAndModify(user, ChangeSetM(setM)); err != nil {
			logger.Error(err)
			abort(ErrInternalError)
			return
		}
	}

	OutputJson(helpers.OutputUserDetailInfo(user), w, r)
	return
}

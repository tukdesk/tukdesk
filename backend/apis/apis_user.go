package apis

import (
	"github.com/labstack/echo"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

var (
	defaultClientsSort = []string{"-updated"}
)

type UserModule struct {
	cfg *config.Config
}

func RegisterUserModule(cfg *config.Config, mux *echo.Group) {
	m := UserModule{
		cfg: cfg,
	}

	group := mux.Group("/users")
	group.Use(CurrentUser)

	group.Get("", m.userList)
	group.Get("/:userId", m.userInfo)
	group.Put("/:userId", m.userUpdate)
	return
}

func (this *UserModule) userList(c *echo.Context) error {
	// 只有 agent 可以查看全员列表
	CheckAuthorizedAsAgent(c)

	r := c.Request()

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

	logger := GetLogger(c)

	count, err := helpers.ClientCount(query)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	users, err := helpers.ClientListAfter(query, listArgs.LastId, listArgs.Limit, sort)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	items := make([]*helpers.OutputUser, len(users))
	if len(items) > 0 {
		for i, user := range users {
			items[i] = helpers.OutputUserProfileInfo(user)
		}
	}

	return OutputJson(ListResultNew(count, items), c)
}

func (this *UserModule) userInfo(c *echo.Context) error {
	userId, ok := helpers.IdFromString(c.Param("userId"))
	if !ok {
		return ErrInvalidId
	}

	logger := GetLogger(c)

	user, err := helpers.UserFindById(userId)
	if helpers.IsNotFound(err) {
		return ErrUserNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrUserNotFound
	}

	var infoParser func(*models.User) *helpers.OutputUser
	if AuthorizedAsAgent(user) {
		infoParser = helpers.OutputUserDetailInfo
	} else {
		infoParser = helpers.OutputUserBaseInfo
	}

	return OutputJson(infoParser(user), c)
}

func (this *UserModule) userUpdate(c *echo.Context) error {
	// 只有 agent 可以通过这个接口修改 user 信息
	// 只能修改 client 信息
	user := CheckAuthorizedAsAgent(c)

	userId, ok := helpers.IdFromString(c.Param("userId"))
	if !ok {
		return ErrInvalidId
	}

	args := GetMapArgsFromContext(c)

	logger := GetLogger(c)

	user, err := helpers.ClientFindById(userId)
	if helpers.IsNotFound(err) {
		return ErrUserNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	setM := helpers.M{}
	v := helpers.ValidationNew()

	for name, val := range args {
		switch name {
		case "business":
			businessMap, ok := val.(map[string]interface{})
			if !ok {
				return ErrorInvalidArgType(name, helpers.JSONTypeNameObject)
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
			return ErrInternalError
		}
	}

	return OutputJson(helpers.OutputUserDetailInfo(user), c)
}

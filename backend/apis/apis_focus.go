package apis

import (
	"github.com/labstack/echo"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

var (
	defaultFocusSort = []string{"-priority", "-deadline", "-created"}
)

type FocusModule struct {
	cfg *config.Config
}

func RegisterFocusModule(cfg *config.Config, mux *echo.Group) {
	m := FocusModule{
		cfg: cfg,
	}

	group := mux.Group("/focus")
	group.Use(CurrentUser)

	group.Get("", m.focusList)
	group.Post("", m.focusAdd)
	group.Put("/:focusId", m.focusHandle)
	return
}

func (this *FocusModule) focusList(c *echo.Context) error {
	CheckAuthorizedAsAgent(c)

	r := c.Request()

	listArgs := GetListArgsFromRequest(r)

	v := helpers.ValidationNew()

	query := helpers.M{}
	if status := r.FormValue("status"); status != "" {
		helpers.ValidationForFocusStatusOnList(v, "status", status)
		CheckValidation(v)
		query["status"] = status
	}

	logger := GetLogger(c)

	count, err := helpers.FocusCount(query)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	items, err := helpers.FocusListAfter(query, listArgs.LastId, listArgs.Limit, defaultFocusSort)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(ListResultNew(count, items), c)
}

func (this *FocusModule) focusAdd(c *echo.Context) error {
	CheckAuthorizedAsAgent(c)

	args := &FocusAddArgs{}
	GetJsonArgsFromContext(c, args)

	v := helpers.ValidationNew()
	helpers.ValidationForFocusResourceType(v, "resourceType", args.ResourceType)
	helpers.ValidationForFocusResourceId(v, "resourceId", args.ResourceId)
	helpers.ValidationForFoucsMessage(v, "message", args.Message)
	helpers.ValidationForFocusPriority(v, "priority", args.Priority)
	CheckValidation(v)

	logger := GetLogger(c)

	resource, err := helpers.ResourceByTypeAndId(args.ResourceType, args.ResourceId)
	if helpers.IsNotFound(err) {
		return ErrResourceNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	focus, err := helpers.FocusReminderInsert(args.Message, args.Duration, args.Priority, resource)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(focus, c)
}

func (this *FocusModule) focusHandle(c *echo.Context) error {
	CheckAuthorizedAsAgent(c)

	focusId, ok := helpers.IdFromString(c.Param("focusId"))
	if !ok {
		return ErrInvalidId
	}

	logger := GetLogger(c)

	focus, err := helpers.FocusFindById(focusId)
	if helpers.IsNotFound(err) {
		return ErrFocusNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	if focus.Status == models.FocusStatusHandled {
		return ErrFocusAlreadyHandled
	}

	if err = helpers.FocusHandled(focus); err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(focus, c)
}

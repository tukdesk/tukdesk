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
	defaultFocusSort = []string{"-priority", "-deadline", "-created"}
)

type FocusModule struct {
	cfg *config.Config
}

func RegisterFocusModule(cfg *config.Config, app *web.Mux) *web.Mux {
	m := FocusModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Get("", m.focusList)
	mux.Post("", m.focusAdd)
	mux.Put("/:focusId", m.focusHandle)
	mux.Use(CurrentUser)

	gojimiddleware.RegisterSubroute("/focus", app, mux)
	return mux
}

func (this *FocusModule) focusList(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c, w, r)

	listArgs := GetListArgsFromRequest(r)

	v := helpers.ValidationNew()

	query := helpers.M{}
	if status := r.FormValue("status"); status != "" {
		helpers.ValidationForFocusStatusOnList(v, "status", status)
		CheckValidation(v)
		query["status"] = status
	}

	logger := GetLogger(&c, w, r)

	count, err := helpers.FocusCount(query)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	items, err := helpers.FocusListAfter(query, listArgs.LastId, listArgs.Limit, defaultFocusSort)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(ListResultNew(count, items), w, r)
	return
}

func (this *FocusModule) focusAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c, w, r)

	args := &FocusAddArgs{}
	GetJsonArgsFromRequest(r, args)

	v := helpers.ValidationNew()
	helpers.ValidationForFocusResourceType(v, "resourceType", args.ResourceType)
	helpers.ValidationForFocusResourceId(v, "resourceId", args.ResourceId)
	helpers.ValidationForFoucsMessage(v, "message", args.Message)
	helpers.ValidationForFocusPriority(v, "priority", args.Priority)
	CheckValidation(v)

	logger := GetLogger(&c, w, r)

	resource, err := helpers.ResourceByTypeAndId(args.ResourceType, args.ResourceId)
	if helpers.IsNotFound(err) {
		abort(err)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	focus, err := helpers.FocusReminderInsert(args.Message, args.Duration, args.Priority, resource)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(focus, w, r)
	return
}

func (this *FocusModule) focusHandle(c web.C, w http.ResponseWriter, r *http.Request) {
	CheckAuthorizedAsAgent(&c, w, r)

	focusId, ok := helpers.IdFromString(c.URLParams["focusId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	logger := GetLogger(&c, w, r)

	focus, err := helpers.FocusFindById(focusId)
	if helpers.IsNotFound(err) {
		abort(ErrFocusNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	if focus.Status == models.FocusStatusHandled {
		abort(ErrFocusAlreadyHandled)
		return
	}

	if err = helpers.FocusHandled(focus); err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(focus, w, r)
	return
}

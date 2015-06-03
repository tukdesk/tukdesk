package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/validation"
	"github.com/zenazn/goji/web"
)

var (
	defaultTicketsSort = []string{"-created"}
)

type TicketModule struct {
	cfg config.Config
}

func RegisterTicketsModule(cfg config.Config, app *web.Mux) *web.Mux {
	m := TicketModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Get("", m.ticketList)
	mux.Post("", m.ticketAdd)

	gojimiddleware.RegisterSubroute("/tickets", app, mux)
	return mux
}

func (this *TicketModule) ticketList(c web.C, w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(&c, w, r)

	isAgent := AuthorizedAsAgent(user)
	isLogged := isAgent || AuthorizedLogged(user)

	listArgs := GetListArgsFromRequest(r)

	v := helpers.ValidationNew()

	if listArgs.Sort != "" {
		v.In("sort", listArgs.Sort, helpers.TicketSortOptionsForList)
	}

	CheckValidation(v)

	filter := helpers.M{}
	if !isAgent {
		if isLogged {
			filter["$or"] = []helpers.M{
				helpers.M{"isPublic": true},
				helpers.M{"creatorId": user.Id},
			}
		} else {
			filter["isPublic"] = true
		}
	}

	filter = helpers.FilterParseFromRequest(r, filter, isAgent)

	sort := make([]string, 0, len(defaultTicketsSort)+1)
	if listArgs.Sort != "" {
		sort = append(sort, listArgs.Sort)
	}
	sort = append(sort, defaultTicketsSort...)

	logger := GetLogger(&c, w, r)

	// get list and count
	count, err := helpers.TicketCount(filter)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	tickets, err := helpers.TicketListAfter(filter, listArgs.LastId, listArgs.Limit, sort)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	items := make([]*helpers.OutputTicketInfo, len(tickets))

	if len(items) > 0 {
		var infoParser func(*models.Ticket) (*helpers.OutputTicketInfo, error)
		if isAgent {
			infoParser = helpers.OutputTicketDetailInfoForList
		} else {
			infoParser = helpers.OutputTicketPublicInfoForList
		}

		for i, ticket := range tickets {
			info, err := infoParser(ticket)
			if err != nil {
				logger.Error(err)
				abort(ErrInternalError)
				return
			}
			items[i] = info
		}
	}

	OutputJson(ListResultNew(count, items), w, r)
	return
}

func (this *TicketModule) ticketAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(&c, w, r)

	args := &TicketAddArgs{}
	GetJsonArgsFromRequest(r, args)

	if args.Subject == "" {
		args.Subject = helpers.TicketGetValidSubject(args.Content)
	}

	if args.Status == "" {
		args.Status = helpers.TicketStatusDefault
	}

	if args.Channel == "" {
		args.Channel = helpers.TicketChannelWeb
	}

	v := helpers.ValidationNew()
	helpers.ValidationForTicketSubject(v, "subject", args.Subject)
	helpers.ValidationForTicketContent(v, "content", args.Content)

	CheckValidation(v)

	var ticketMaker func(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *gojimiddleware.XLogger) *models.Ticket
	if AuthorizedAsAgent(user) {
		ticketMaker = this.ticketMakerForAgent
	} else if AuthorizedLogged(user) {
		ticketMaker = this.ticketMakerForClient
	} else {
		ticketMaker = this.ticketMakerForAnonym
	}

	logger := GetLogger(&c, w, r)

	ticket := ticketMaker(user, args, v, logger)

	// check automation

	err := helpers.TicketInit(ticket)
	if helpers.IsDup(err) {
		abort(ErrTicketDuplicate)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	info, err := helpers.OutputTicketPublicInfo(ticket)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(info, w, r)
	return
}

func (this *TicketModule) ticketMakerForAnonym(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *gojimiddleware.XLogger) *models.Ticket {
	// 需要 email; 不可设置 status, isPublic; extend 受限
	extend := helpers.TicketParseExtendFromPreSet(args.Extend)

	helpers.ValidationForEmail(v, "email", args.Email)
	helpers.ValidationForTicketExtendField(v, extend)

	CheckValidation(v)

	creator, _, err := helpers.UserMustForChannelEmail(args.Email, helpers.UserGetValidNameFromEmail(args.Email))
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return nil
	}

	return helpers.TicketNewWithChannelName(creator, args.Channel, args.Subject, args.Content, extend)
}

func (this *TicketModule) ticketMakerForClient(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *gojimiddleware.XLogger) *models.Ticket {
	// 不需要 email; 不可设置 status; extend 受限
	extend := helpers.TicketParseExtendFromPreSet(args.Extend)

	helpers.ValidationForTicketExtendField(v, extend)

	CheckValidation(v)

	ticket := helpers.TicketNewWithChannelName(user, args.Channel, args.Subject, args.Content, extend)
	ticket.IsPublic = args.IsPublic
	return ticket
}

func (this *TicketModule) ticketMakerForAgent(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *gojimiddleware.XLogger) *models.Ticket {
	// 需要 email; extend 不受限

	helpers.ValidationForEmail(v, "email", args.Email)
	helpers.ValidationForTicketExtendField(v, args.Extend)
	helpers.ValidationForTicketStatusOnCreate(v, "status", args.Status)

	CheckValidation(v)

	creator, _, err := helpers.UserMustForChannelEmail(args.Email, helpers.UserGetValidName(args.Email))
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return nil
	}

	// todo accept chId for agent?
	ticket := helpers.TicketNewWithChannelName(creator, args.Channel, args.Subject, args.Content, args.Extend)
	ticket.IsPublic = args.IsPublic
	ticket.Status = args.Status

	return ticket
}

func (this *TicketModule) ticketProfile(c web.C, w http.ResponseWriter, r *http.Request) {

	return
}

func (this *TicketModule) ticketUpdate(c web.C, w http.ResponseWriter, r *http.Request) {

	return
}

func (this *TicketModule) commentList(c web.C, w http.ResponseWriter, r *http.Request) {

	return
}

func (this *TicketModule) commentAdd(c web.C, w http.ResponseWriter, r *http.Request) {

	return
}

func (this *TicketModule) commentUpdate(c web.C, w http.ResponseWriter, r *http.Request) {

	return
}

package apis

import (
	"net/http"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/validation"
	"github.com/tukdesk/httputils/xlogger"
	"github.com/zenazn/goji/web"
)

var (
	defaultTicketsSort = []string{"-created"}
	defaultCommentSort = []string{"updated"}
)

type TicketModule struct {
	cfg *config.Config
}

func RegisterTicketsModule(cfg *config.Config, app *web.Mux) *web.Mux {
	m := TicketModule{
		cfg: cfg,
	}

	mux := web.New()
	mux.Get("", m.ticketList)
	mux.Post("", m.ticketAdd)
	mux.Get("/:ticketId", m.ticketInfo)
	mux.Put("/:ticketId", m.ticketUpdate)
	mux.Get("/:ticketId/comments", m.commentList)
	mux.Post("/:ticketId/comments", m.commentAdd)
	mux.Put("/:ticketId/comments/:commentId", m.commentUpdate)
	mux.Use(CurrentUser)

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
		helpers.ValidationForTicektListSort(v, "sort", listArgs.Sort)
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

	items := make([]*helpers.OutputTicket, len(tickets))

	if len(items) > 0 {
		var infoParser func(*models.Ticket) (*helpers.OutputTicket, error)
		if isAgent {
			infoParser = helpers.OutputTicketDetailInfoForList
		} else {
			infoParser = helpers.OutputTicketPublicInfoForList
		}

		commentQuery := helpers.M{}
		showComments := r.FormValue("nocomments") != trueInQuery
		if showComments {
			if !AuthorizedAsAgent(user) {
				commentQuery["type"] = helpers.M{"$in": helpers.CommentTypeOptionsForNonAgentView}
			}
		}

		for i, ticket := range tickets {
			info, err := infoParser(ticket)
			if err != nil {
				logger.Error(err)
				abort(ErrInternalError)
				return
			}
			if showComments {
				if err = info.GetComments(commentQuery, defaultCommentSort); err != nil {
					logger.Error(err)
					abort(ErrInternalError)
					return
				}
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

	args.Extend = helpers.TicketParseExtendFromPreSet(args.Extend)

	v := helpers.ValidationNew()
	helpers.ValidationForTicketSubject(v, "subject", args.Subject)
	helpers.ValidationForTicketContent(v, "content", args.Content)
	helpers.ValidationForTicketChannel(v, "channel", args.Channel)
	helpers.ValidationForTicketExtendField(v, args.Extend)

	CheckValidation(v)

	var ticketMaker func(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *xlogger.XLogger) *models.Ticket
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

	err := helpers.TicketInit(ticket, args.Content, args.Attachments)
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

func (this *TicketModule) ticketMakerForAnonym(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *xlogger.XLogger) *models.Ticket {
	// 需要 email; 不可设置 status, isPublic;
	helpers.ValidationForEmail(v, "email", args.Email)

	CheckValidation(v)

	creator, _, err := helpers.UserMustForChannelEmail(args.Email, helpers.UserGetValidNameFromEmail(args.Email))
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return nil
	}

	return helpers.TicketNewWithChannelName(creator, args.Channel, args.Subject, args.Extend)
}

func (this *TicketModule) ticketMakerForClient(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *xlogger.XLogger) *models.Ticket {
	// 不需要 email; 不可设置 status;
	ticket := helpers.TicketNewWithChannelName(user, args.Channel, args.Subject, args.Extend)
	ticket.IsPublic = args.IsPublic
	return ticket
}

func (this *TicketModule) ticketMakerForAgent(user *models.User, args *TicketAddArgs, v *validation.Validation, logger *xlogger.XLogger) *models.Ticket {
	// 需要 email;
	helpers.ValidationForEmail(v, "email", args.Email)
	helpers.ValidationForTicketStatusOnCreate(v, "status", args.Status)

	CheckValidation(v)

	creator, _, err := helpers.UserMustForChannelEmail(args.Email, helpers.UserGetValidNameFromEmail(args.Email))
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return nil
	}

	// todo accept chId for agent?
	ticket := helpers.TicketNewWithChannelName(creator, args.Channel, args.Subject, args.Extend)
	ticket.IsPublic = args.IsPublic
	ticket.Status = args.Status

	return ticket
}

func (this *TicketModule) ticketInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(&c, w, r)

	ticketId, ok := helpers.IdFromString(c.URLParams["ticketId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	logger := GetLogger(&c, w, r)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		abort(ErrTicketNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	showDetail := AuthorizedAsAgent(user) || AuthorizedAsSpecifiedUser(user, ticket.CreatorId)
	var infoParser func(*models.Ticket) (*helpers.OutputTicket, error)
	if showDetail {
		infoParser = helpers.OutputTicketDetailInfo
	} else {
		infoParser = helpers.OutputTicketPublicInfo
	}

	output, err := infoParser(ticket)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(output, w, r)
	return
}

func (this *TicketModule) ticketUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// 仅修改工单属性, 不修改内容, 因此不向 client 开放
	CheckAuthorizedAsAgent(&c, w, r)

	ticketId, ok := helpers.IdFromString(c.URLParams["ticketId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	args := GetMapArgsFromRequest(r)

	logger := GetLogger(&c, w, r)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		abort(ErrTicketNotFound)
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
		case "priority":
			priorityStr, ok := val.(string)
			if !ok {
				abort(ErrorInvalidArgType(name, helpers.JSONTypeNameString))
				return
			}

			priority := models.NewTypePriority(priorityStr)

			// 不先比较, 因为需要将输出值转化为真实值
			if priority == ticket.Priority {
				break
			}

			helpers.ValidationForTicketPriority(v, name, priority)
			CheckValidation(v)

			setM[name] = priority

		case "isPublic":
			if val == ticket.IsPublic {
				break
			}

			isPublic, ok := val.(bool)
			if !ok {
				abort(ErrorInvalidArgType(name, helpers.JSONTypeNameBoolen))
				return
			}

			setM[name] = isPublic

		case "status":
			if val == ticket.Status {
				break
			}

			status, ok := val.(string)
			if !ok {
				abort(ErrorInvalidArgType(name, helpers.JSONTypeNameString))
				return
			}

			helpers.ValidationForTicketStatusOnUpdate(v, name, status)
			CheckValidation(v)

			setM[name] = val

		}
	}

	if len(setM) > 0 {
		setM["updated"] = NowUnix()

		err := helpers.TicketFindAndModify(ticket, ChangeSetM(setM))
		if helpers.IsNotFound(err) {
			abort(ErrTicketNotFound)
			return
		}

		if err != nil {
			logger.Error(err)
			abort(ErrInternalError)
			return
		}
	}

	output, err := helpers.OutputTicketDetailInfo(ticket)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(output, w, r)
	return
}

func (this *TicketModule) commentList(c web.C, w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(&c, w, r)

	ticketId, ok := helpers.IdFromString(c.URLParams["ticketId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	logger := GetLogger(&c, w, r)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		abort(ErrTicketNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	// 非公开 ticket, agent 或 题主可见
	if !ticket.IsPublic && !AuthorizedAsAgent(user) && !AuthorizedAsSpecifiedUser(user, ticket.CreatorId) {
		abort(ErrUnauthorized)
		return
	}

	// 游客可见: public, question, feedback
	// 题主可见: public, question, feedback
	// agent 可见: all
	query := helpers.M{}
	if !AuthorizedAsAgent(user) {
		query["type"] = helpers.M{"$in": helpers.CommentTypeOptionsForNonAgentView}
	}

	comments, err := helpers.CommentFindAllByTicketId(ticket.Id, query, defaultCommentSort)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	items := make([]*helpers.OutputComment, len(comments))

	if len(items) > 0 {
		for i, comment := range comments {
			info, err := helpers.OutputCommentInfo(comment)
			if err != nil {
				logger.Error(err)
				abort(ErrInternalError)
				return
			}
			items[i] = info
		}
	}

	OutputJson(ListResultNew(len(items), items), w, r)
	return
}

func (this *TicketModule) commentAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	user := CheckAuthorizedLogged(&c, w, r)

	ticketId, ok := helpers.IdFromString(c.URLParams["ticketId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	args := &CommentAddArgs{}
	GetJsonArgsFromRequest(r, args)

	v := helpers.ValidationNew()
	helpers.ValidationForCommentTypeOnCreate(v, "type", args.Type)
	helpers.ValidationForCommentContent(v, "content", args.Content)
	CheckValidation(v)

	switch args.Type {
	case models.CommentTypePublic, models.CommentTypeInternal:
		if !AuthorizedAsAgent(user) {
			abort(ErrUnauthorized)
			return
		}
	}

	logger := GetLogger(&c, w, r)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		abort(ErrTicketNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	if !AuthorizedAsAgent(user) && !AuthorizedAsSpecifiedUser(user, ticket.CreatorId) {
		abort(ErrUnauthorized)
		return
	}

	comment, err := helpers.CommentInsertForTicket(ticket, user, args.Type, args.Content, args.Attachments)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	ticketSetM := helpers.M{}
	if args.Type == models.CommentTypeFeedback {
		ticketSetM["status"] = models.TicketStatusResubmitted
	} else if args.Type == models.CommentTypePublic {
		ticketSetM["status"] = models.TicketStatusReplied
		if ticket.FirstCommented == 0 {
			ticketSetM["firstCommented"] = NowUnix()
		}
	}

	if len(ticketSetM) > 0 {
		ticketSetM["updated"] = NowUnix()
		if err := helpers.TicketFindAndModify(ticket, ChangeSetM(ticketSetM)); err != nil {
			logger.Error(err)
			abort(ErrInternalError)
			return
		}
	}

	output, err := helpers.OutputCommentInfo(comment)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(output, w, r)
	return
}

func (this *TicketModule) commentUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	// 当 comment 类型为 internal 时, 可以修改内容, 或将类型修改为 public
	CheckAuthorizedAsAgent(&c, w, r)

	ticketId, ok := helpers.IdFromString(c.URLParams["ticketId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	commentId, ok := helpers.IdFromString(c.URLParams["commentId"])
	if !ok {
		abort(ErrInvalidId)
		return
	}

	args := GetMapArgsFromRequest(r)

	// get ticket and comment
	logger := GetLogger(&c, w, r)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		abort(ErrTicketNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	comment, err := helpers.CommentFindByTicketIdAndCommentId(ticketId, commentId)
	if helpers.IsNotFound(err) {
		abort(ErrCommentNotFound)
		return
	}

	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	// 非 internal 的 comment, 不可修改
	if comment.Type != models.CommentTypeInternal {
		abort(ErrCommentUnchangeable)
		return
	}

	v := helpers.ValidationNew()

	commentSetM := helpers.M{}
	ticketSetM := helpers.M{}

	for name, val := range args {
		switch name {
		case "type":
			if val == comment.Type {
				break
			}

			commentType, ok := val.(string)
			if !ok {
				abort(ErrorInvalidArgType(name, helpers.JSONTypeNameString))
				return
			}

			helpers.ValidationForCommentTypeOnUpdate(v, name, commentType)
			CheckValidation(v)

			commentSetM[name] = commentType
			// 假装我们不知道 Update 时 只允许 CommentTypePublic 一种
			if commentType == models.CommentTypePublic {
				ticketSetM["status"] = models.TicketStatusReplied
			}

		case "content":
			if val == comment.Content {
				break
			}

			content, ok := val.(string)
			if !ok {
				abort(ErrorInvalidArgType(name, helpers.JSONTypeNameString))
				return
			}

			commentSetM[name] = content
		}
	}

	now := NowUnix()

	if len(ticketSetM) > 0 {
		ticketSetM["updated"] = now
		if err := helpers.TicketFindAndModify(ticket, ChangeSetM(ticketSetM)); err != nil {
			logger.Error(err)
			abort(ErrInternalError)
			return
		}
	}

	if len(commentSetM) > 0 {
		commentSetM["updated"] = now
		if err := helpers.CommentFindAndModify(comment, ChangeSetM(commentSetM)); err != nil {
			logger.Error(err)
			abort(ErrInternalError)
			return
		}
	}

	output, err := helpers.OutputCommentInfo(comment)
	if err != nil {
		logger.Error(err)
		abort(ErrInternalError)
		return
	}

	OutputJson(output, w, r)
	return
}

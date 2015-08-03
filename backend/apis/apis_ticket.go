package apis

import (
	"github.com/labstack/echo"
	"github.com/tukdesk/httputils/validation"

	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models"
	"github.com/tukdesk/tukdesk/backend/models/helpers"
)

var (
	defaultTicketsSort = []string{"-created"}
	defaultCommentSort = []string{"updated"}
)

type TicketModule struct {
	cfg *config.Config
}

func RegisterTicketsModule(cfg *config.Config, mux *echo.Group) {
	m := TicketModule{
		cfg: cfg,
	}

	group := mux.Group("/tickets")
	group.Use(CurrentUser)

	group.Get("", m.ticketList)
	group.Post("", m.ticketAdd)
	group.Get("/:ticketId", m.ticketInfo)
	group.Put("/:ticketId", m.ticketUpdate)
	group.Get("/:ticketId/comments", m.commentList)
	group.Post("/:ticketId/comments", m.commentAdd)
	group.Put("/:ticketId/comments/:commentId", m.commentUpdate)
	return
}

func (this *TicketModule) ticketList(c *echo.Context) error {
	user := GetCurrentUser(c)

	r := c.Request()

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

	logger := GetLogger(c)

	// get list and count
	count, err := helpers.TicketCount(filter)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	tickets, err := helpers.TicketListAfter(filter, listArgs.LastId, listArgs.Limit, sort)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
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
				return ErrInternalError
			}
			if showComments {
				if err = info.GetComments(commentQuery, defaultCommentSort); err != nil {
					logger.Error(err)
					return ErrInternalError
				}
			}
			items[i] = info
		}
	}

	return OutputJson(ListResultNew(count, items), c)
}

func (this *TicketModule) ticketAdd(c *echo.Context) error {
	user := GetCurrentUser(c)

	args := &TicketAddArgs{}
	GetJsonArgsFromContext(c, args)

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

	var ticketMaker func(user *models.User, args *TicketAddArgs, v *validation.Validation) (*models.Ticket, error)
	if AuthorizedAsAgent(user) {
		ticketMaker = this.ticketMakerForAgent
	} else if AuthorizedLogged(user) {
		ticketMaker = this.ticketMakerForClient
	} else {
		ticketMaker = this.ticketMakerForAnonym
	}

	logger := GetLogger(c)

	ticket, err := ticketMaker(user, args, v)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	// check automation

	err = helpers.TicketInit(ticket, args.Content, args.Attachments)
	if helpers.IsDup(err) {
		return ErrTicketDuplicate
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	info, err := helpers.OutputTicketPublicInfo(ticket)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(info, c)
}

func (this *TicketModule) ticketMakerForAnonym(user *models.User, args *TicketAddArgs, v *validation.Validation) (*models.Ticket, error) {
	// 需要 email; 不可设置 status, isPublic;
	helpers.ValidationForEmail(v, "email", args.Email)

	CheckValidation(v)

	creator, _, err := helpers.UserMustForChannelEmail(args.Email, helpers.UserGetValidNameFromEmail(args.Email))
	if err != nil {
		return nil, err
	}

	return helpers.TicketNewWithChannelName(creator, args.Channel, args.Subject, args.Extend), nil
}

func (this *TicketModule) ticketMakerForClient(user *models.User, args *TicketAddArgs, v *validation.Validation) (*models.Ticket, error) {
	// 不需要 email; 不可设置 status;
	ticket := helpers.TicketNewWithChannelName(user, args.Channel, args.Subject, args.Extend)
	ticket.IsPublic = args.IsPublic
	return ticket, nil
}

func (this *TicketModule) ticketMakerForAgent(user *models.User, args *TicketAddArgs, v *validation.Validation) (*models.Ticket, error) {
	// 需要 email;
	helpers.ValidationForEmail(v, "email", args.Email)
	helpers.ValidationForTicketStatusOnCreate(v, "status", args.Status)

	CheckValidation(v)

	creator, _, err := helpers.UserMustForChannelEmail(args.Email, helpers.UserGetValidNameFromEmail(args.Email))
	if err != nil {
		return nil, err
	}

	// todo accept chId for agent?
	ticket := helpers.TicketNewWithChannelName(creator, args.Channel, args.Subject, args.Extend)
	ticket.IsPublic = args.IsPublic
	ticket.Status = args.Status

	return ticket, nil
}

func (this *TicketModule) ticketInfo(c *echo.Context) error {
	user := GetCurrentUser(c)

	ticketId, ok := helpers.IdFromString(c.Param("ticketId"))
	if !ok {
		return ErrInvalidId
	}

	logger := GetLogger(c)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		return ErrTicketNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
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
		return ErrInternalError
	}

	return OutputJson(output, c)
}

func (this *TicketModule) ticketUpdate(c *echo.Context) error {
	// 仅修改工单属性, 不修改内容, 因此不向 client 开放
	CheckAuthorizedAsAgent(c)

	ticketId, ok := helpers.IdFromString(c.Param("ticketId"))
	if !ok {
		return ErrInvalidId
	}

	args := GetMapArgsFromContext(c)

	logger := GetLogger(c)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		return ErrTicketNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	setM := helpers.M{}
	v := helpers.ValidationNew()

	for name, val := range args {
		switch name {
		case "priority":
			priorityStr, ok := val.(string)
			if !ok {
				return ErrorInvalidArgType(name, helpers.JSONTypeNameString)
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
				return ErrorInvalidArgType(name, helpers.JSONTypeNameBoolen)
			}

			setM[name] = isPublic

		case "status":
			if val == ticket.Status {
				break
			}

			status, ok := val.(string)
			if !ok {
				return ErrorInvalidArgType(name, helpers.JSONTypeNameString)
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
			return ErrTicketNotFound
		}

		if err != nil {
			logger.Error(err)
			return ErrInternalError
		}
	}

	output, err := helpers.OutputTicketDetailInfo(ticket)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(output, c)
}

func (this *TicketModule) commentList(c *echo.Context) error {
	user := GetCurrentUser(c)

	ticketId, ok := helpers.IdFromString(c.Param("ticketId"))
	if !ok {
		return ErrInvalidId
	}

	logger := GetLogger(c)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		return ErrTicketNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	// 非公开 ticket, agent 或 题主可见
	if !ticket.IsPublic && !AuthorizedAsAgent(user) && !AuthorizedAsSpecifiedUser(user, ticket.CreatorId) {
		return ErrUnauthorized
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
		return ErrInternalError
	}

	items := make([]*helpers.OutputComment, len(comments))

	if len(items) > 0 {
		for i, comment := range comments {
			info, err := helpers.OutputCommentInfo(comment)
			if err != nil {
				logger.Error(err)
				return ErrInternalError
			}
			items[i] = info
		}
	}

	return OutputJson(ListResultNew(len(items), items), c)
}

func (this *TicketModule) commentAdd(c *echo.Context) error {
	user := CheckAuthorizedLogged(c)

	ticketId, ok := helpers.IdFromString(c.Param("ticketId"))
	if !ok {
		return ErrInvalidId
	}

	args := &CommentAddArgs{}
	GetJsonArgsFromContext(c, args)

	v := helpers.ValidationNew()
	helpers.ValidationForCommentTypeOnCreate(v, "type", args.Type)
	helpers.ValidationForCommentContent(v, "content", args.Content)
	CheckValidation(v)

	switch args.Type {
	case models.CommentTypePublic, models.CommentTypeInternal:
		if !AuthorizedAsAgent(user) {
			return ErrUnauthorized
		}
	}

	logger := GetLogger(c)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		return ErrTicketNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	if !AuthorizedAsAgent(user) && !AuthorizedAsSpecifiedUser(user, ticket.CreatorId) {
		return ErrUnauthorized
	}

	comment, err := helpers.CommentInsertForTicket(ticket, user, args.Type, args.Content, args.Attachments)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
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
			return ErrInternalError
		}
	}

	output, err := helpers.OutputCommentInfo(comment)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(output, c)
}

func (this *TicketModule) commentUpdate(c *echo.Context) error {
	// 当 comment 类型为 internal 时, 可以修改内容, 或将类型修改为 public
	CheckAuthorizedAsAgent(c)

	ticketId, ok := helpers.IdFromString(c.Param("ticketId"))
	if !ok {
		return ErrInvalidId
	}

	commentId, ok := helpers.IdFromString(c.Param("commentId"))
	if !ok {
		return ErrInvalidId
	}

	args := GetMapArgsFromContext(c)

	// get ticket and comment
	logger := GetLogger(c)

	ticket, err := helpers.TicketFindById(ticketId)
	if helpers.IsNotFound(err) {
		return ErrTicketNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	comment, err := helpers.CommentFindByTicketIdAndCommentId(ticketId, commentId)
	if helpers.IsNotFound(err) {
		return ErrCommentNotFound
	}

	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	// 非 internal 的 comment, 不可修改
	if comment.Type != models.CommentTypeInternal {
		return ErrCommentUnchangeable
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
				return ErrorInvalidArgType(name, helpers.JSONTypeNameString)
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
				return ErrorInvalidArgType(name, helpers.JSONTypeNameString)
			}

			commentSetM[name] = content
		}
	}

	now := NowUnix()

	if len(ticketSetM) > 0 {
		ticketSetM["updated"] = now
		if err := helpers.TicketFindAndModify(ticket, ChangeSetM(ticketSetM)); err != nil {
			logger.Error(err)
			return ErrInternalError
		}
	}

	if len(commentSetM) > 0 {
		commentSetM["updated"] = now
		if err := helpers.CommentFindAndModify(comment, ChangeSetM(commentSetM)); err != nil {
			logger.Error(err)
			return ErrInternalError
		}
	}

	output, err := helpers.OutputCommentInfo(comment)
	if err != nil {
		logger.Error(err)
		return ErrInternalError
	}

	return OutputJson(output, c)
}

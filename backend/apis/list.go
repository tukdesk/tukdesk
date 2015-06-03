package apis

import (
	"net/http"
	"strconv"

	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"gopkg.in/mgo.v2/bson"
)

const (
	ListLimitMin     = 5
	ListLimitDefault = 30
	ListLimitMax     = 100
)

type ListArgs struct {
	Limit  int           `json:"limit"`
	Sort   string        `json:"sort"`
	LastId bson.ObjectId `json:"lastId"`
}

func GetListArgsFromRequest(r *http.Request) *ListArgs {
	args := &ListArgs{}

	args.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))

	if args.Limit <= 0 {
		args.Limit = ListLimitDefault
	} else if args.Limit < ListLimitMin {
		args.Limit = ListLimitMin
	} else if args.Limit > ListLimitMax {
		args.Limit = ListLimitMax
	}

	args.Sort = r.URL.Query().Get("sort")

	args.LastId = helpers.EmptyId

	lastIdStr := r.URL.Query().Get("lastId")
	if lastIdStr != "" {
		if lastId, ok := helpers.IdFromString(lastIdStr); ok {
			args.LastId = lastId
		}
	}

	return args
}

type ListResult struct {
	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

func ListResultNew(count int, items interface{}) *ListResult {
	return &ListResult{
		Count: count,
		Items: items,
	}
}

package helpers

import (
	"net/http"
	"strings"
	"time"

	"github.com/tukdesk/tukdesk/backend/models"

	"github.com/tukdesk/httputils/validation"
)

var (
	FilterFieldValueOptions = map[string][]interface{}{
		"status":   TicketStatusOptionsForList,
		"priority": TicketPriorityOptions,
	}

	FilterFieldNameOptions = []interface{}{
		"creatorId",
		"priority",
		"status",
		"created",
		"updated",
		"sort",
	}

	FilterTimeFieldOptions = []interface{}{
		"1day",   // 一天
		"2days",  // 两天
		"3days",  // 三天
		"7days",  // 一周
		"14days", // 两周
		"30days", // 30天
	}

	oneDay = 24 * time.Hour

	FilterTimeFieldDurationMap = map[string]time.Duration{
		"1day":   oneDay,
		"2days":  2 * oneDay,
		"3days":  3 * oneDay,
		"7days":  7 * oneDay,
		"14days": 14 * oneDay,
		"30days": 30 * oneDay,
	}
)

func FilterParseFromRequest(req *http.Request, query map[string]interface{}, isAgent bool) map[string]interface{} {
	if query == nil {
		query = M{}
	}

	for _, opt := range FilterFieldNameOptions {
		str := opt.(string)
		valueStr := req.FormValue(str)
		if valueStr == "" {
			continue
		}

		switch str {
		case "creatorId":
			if creatorId, ok := IdFromString(valueStr); ok {
				query[str] = creatorId
			}
		case "updated", "created":
			if duration, ok := FilterTimeFieldDurationMap[valueStr]; ok {
				then := time.Now().Add(-duration).Unix()
				query[str] = M{"$gt": then}
			}

		case "status":
			pieces := strings.Split(valueStr, ",")
			options := FilterFieldValueOptions[str]

			statusValues := []string{}

			exists := map[string]int8{}

			if options != nil {
				for _, piece := range pieces {
					if validation.ValidatorIn(piece, options) && exists[piece] == 0 {
						statusValues = append(statusValues, piece)
						exists[piece] = 1
					}
				}

				if len(statusValues) == 1 {
					query[str] = statusValues[0]
				} else if len(statusValues) > 1 {
					query[str] = M{"$in": statusValues}
				}
			}

		case "priority":
			if !isAgent {
				continue
			}

			strPieces := strings.Split(valueStr, ",")
			options := FilterFieldValueOptions[str]

			priorityValues := []models.TypePriority{}
			exists := map[models.TypePriority]int8{}

			if options != nil {
				for _, strPiece := range strPieces {
					priority := models.NewTypePriority(strPiece)
					if validation.ValidatorIn(priority, options) && exists[priority] == 0 {
						priorityValues = append(priorityValues, priority)
						exists[priority] = 1
					}
				}

				if len(priorityValues) == 1 {
					query[str] = priorityValues[0]
				} else if len(priorityValues) > 1 {
					query[str] = M{"$in": priorityValues}
				}
			}
		}
	}

	return query
}

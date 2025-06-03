package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	SUCCESS = 1
	FAIL    = 0
)

func ResponseSuccess(data interface{}) fiber.Map {
	t := time.Now()
	return fiber.Map{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    SUCCESS,
		"items":     data,
		"error":     nil,
	}
}

/*
Return list data with pagination infos
*/
// NOTE: if not pagination infos is not required, pass -1 to currentPage, currentPageTotalItem, totalPage
func ResponseSuccessListData(data interface{}, currentPage, limit, TotalItem, totalPage int) fiber.Map {
	t := time.Now()
	return fiber.Map{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    SUCCESS,
		"items": fiber.Map{
			"list_data": data,
			"pagination": fiber.Map{
				"current_page": currentPage,
				"total_item":   TotalItem,
				"total_page":   totalPage,
				"limit":        limit,
			},
		},
		"error": nil,
	}
}

/*
Return list data with pagination and static of totalcount that system query count
*/
// NOTE: if not pagination infos is not required, pass -1 to currentPage, currentPageTotalItem, totalPage, totalCount
func ResponseSuccessListDataV2(data interface{}, currentPage, currentPageTotalItem, totalPage int, totalCount int) fiber.Map {
	t := time.Now()
	return fiber.Map{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    SUCCESS,
		"items": fiber.Map{
			"list_data": data,
			"pagination": fiber.Map{
				"current_page":            currentPage,
				"current_page_total_item": currentPageTotalItem,
				"total_page":              totalPage,
				"total_count":             totalCount,
			},
		},
		"error": nil,
	}
}

/*
Return list data with pagination infos, and extra data.

If not pagination infos is not required, pass -1 to currentPage, currentPageTotalItem, totalPage.
*/
func ResponseSuccessListDataWithExtras(data interface{}, currentPage, currentPageTotalItem, totalPage int, extraData interface{}) fiber.Map {
	t := time.Now()
	return fiber.Map{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    SUCCESS,
		"items": fiber.Map{
			"list_data":  data,
			"extra_data": extraData,
			"pagination": fiber.Map{
				"current_page":            currentPage,
				"current_page_total_item": currentPageTotalItem,
				"total_page":              totalPage,
			},
		},
		"error": nil,
	}
}

func ResponseError(data interface{}) fiber.Map {
	t := time.Now()
	return fiber.Map{
		"timestamp": t.Format("2006-01-02-15-04-05"),
		"status":    FAIL,
		"items":     nil,
		"error":     data,
	}
}

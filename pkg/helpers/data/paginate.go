package data

import (
	"crowdfund/pkg/helpers/convertor"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	PageSize   int         `json:"-"`
	PageNumber int         `json:"-"`
	Offset     int         `json:"-"`
	TotalPage  int64       `json:"total_page"`
	TotalRow   int64       `json:"total_row"`
	Items      interface{} `json:"items"`
}

func Paginate(c *fiber.Ctx, totalRows int64) *Pagination {
	pageSize := convertor.StringToInt(c.Query("page_size", "50"))
	pageNumber := convertor.StringToInt(c.Query("page_number", "1"))
	offset := pageSize * (pageNumber - 1)
	pagination := &Pagination{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		TotalRow:   totalRows,
		Offset:     offset,
		TotalPage:  (totalRows + int64(pageSize) - 1) / int64(pageSize),
	}
	return pagination
}

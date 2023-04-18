package organization

import (
	"crowdfund/pkg/core"

	"github.com/gofiber/fiber/v2"
)

type OrgTypeHandler struct{}

func (*OrgTypeHandler) OrgTypeList(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = OrgTypeList(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

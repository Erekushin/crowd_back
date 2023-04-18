package user

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/convertor"

	"github.com/gofiber/fiber/v2"
)

type UserOrgHandler struct{}

func (*UserOrgHandler) OrgList(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)
	userId := convertor.StringToInt(c.Query("user_id"))
	if userId == 0 {
		return core.Resolve(400, c, core.Response("user_id is required"))
	}

	if res, err = OrgList(userId); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

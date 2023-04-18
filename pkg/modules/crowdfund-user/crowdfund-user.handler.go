package crowdfunduser

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CrowdfundUserHandler struct{}

func (*CrowdfundUserHandler) Create(c *fiber.Ctx) error {
	req := new(ReqCrowdfundUserCreate)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := CreateCrowdfundUser(req, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CrowdfundUserHandler) ListUser(c *fiber.Ctx) error {

	var (
		res interface{}
		err error
	)
	CrowdfundID := convertor.StringToInt(c.Query("crowdfund_id"))

	if res, err = CrowdfundUserList(CrowdfundID, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*CrowdfundUserHandler) ListUserCrowdfund(c *fiber.Ctx) error {

	var (
		res interface{}
		err error
	)
	UserId := oauth.GetSessionUserId(c)
	fmt.Println("UserId", UserId)
	if res, err = UserCrowdfundList(UserId, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

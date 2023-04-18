package crowdfund

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CrowdfundHandler struct{}

func (*CrowdfundHandler) List(c *fiber.Ctx) error {

	var (
		res interface{}
		err error
	)

	if res, err = List(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*CrowdfundHandler) Create(c *fiber.Ctx) error {
	req := new(ReqCrowdfundCreate)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := CreateCrowndfund(req, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CrowdfundHandler) Update(c *fiber.Ctx) error {
	crowdfund := new(Crowdfund)

	if err := c.BodyParser(crowdfund); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*crowdfund); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}
	crowdfund.UpdatedBy = oauth.GetSessionUserId(c)
	if err := crowdfund.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CrowdfundHandler) Delete(c *fiber.Ctx) error {
	crowdfund := new(Crowdfund)

	if err := c.BodyParser(crowdfund); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*crowdfund); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := crowdfund.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CrowdfundHandler) ListConfirmed(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = ListConfirmed(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))

}

func (*CrowdfundHandler) Confirm(c *fiber.Ctx) error {
	req := new(ReqCrowdfundId)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := ConfirmCrowndfund(req, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CrowdfundHandler) Cancel(c *fiber.Ctx) error {
	req := new(ReqCrowdfundId)

	if err := c.BodyParser(req); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := CancelCrowdfund(req, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

type Response struct {
	Code   uint        `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func (*CrowdfundHandler) Info(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)
	req := new(ReqCrowdfundId)
	if err := c.BodyParser(req); err != nil {
		fmt.Println("sd")
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*req); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if res, err = Info(req, c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

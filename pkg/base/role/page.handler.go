package role

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type PageHandler struct{}

func (*PageHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = PageList(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*PageHandler) Create(c *fiber.Ctx) error {
	page := new(Page)

	if err := c.BodyParser(page); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*page); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	page.CreatedBy = oauth.GetSessionUserId(c)
	if err := page.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*PageHandler) Update(c *fiber.Ctx) error {
	page := new(Page)

	if err := c.BodyParser(page); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*page); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	page.UpdatedBy = oauth.GetSessionUserId(c)
	if err := page.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*PageHandler) Delete(c *fiber.Ctx) error {
	page := new(Page)

	if err := c.BodyParser(page); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*page); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	page.DeletedBy = oauth.GetSessionUserId(c)
	if err := page.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

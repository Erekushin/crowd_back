package news

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct{}

func (*NewsHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	orgId := oauth.GetSessionOrgId(c)
	if res, err = NewsList(c, orgId); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*NewsHandler) Create(c *fiber.Ctx) error {
	page := new(News)

	if err := c.BodyParser(page); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*page); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}
	page.OrgId = oauth.GetSessionOrgId(c)
	page.CreatedBy = oauth.GetSessionUserId(c)
	if err := page.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*NewsHandler) Update(c *fiber.Ctx) error {
	page := new(News)

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

func (*NewsHandler) Delete(c *fiber.Ctx) error {
	page := new(News)

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

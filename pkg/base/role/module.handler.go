package role

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type ModuleHandler struct{}

func (*ModuleHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = ModuleList(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*ModuleHandler) Create(c *fiber.Ctx) error {
	module := new(Module)

	if err := c.BodyParser(module); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*module); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	module.CreatedBy = oauth.GetSessionUserId(c)
	if err := module.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*ModuleHandler) Update(c *fiber.Ctx) error {
	module := new(Module)

	if err := c.BodyParser(module); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*module); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	module.UpdatedBy = oauth.GetSessionUserId(c)
	if err := module.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*ModuleHandler) Delete(c *fiber.Ctx) error {
	module := new(Module)

	if err := c.BodyParser(module); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*module); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	module.DeletedBy = oauth.GetSessionUserId(c)
	if err := module.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

package document

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type TypeHandler struct{}

func (*TypeHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = TypeList(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*TypeHandler) Create(c *fiber.Ctx) error {
	data := new(Type)
	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.CreatedBy = oauth.GetSessionUserId(c)
	if err := data.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*TypeHandler) Update(c *fiber.Ctx) error {
	data := new(Type)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.UpdatedBy = oauth.GetSessionUserId(c)
	if err := data.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*TypeHandler) Delete(c *fiber.Ctx) error {
	data := new(Type)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := data.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

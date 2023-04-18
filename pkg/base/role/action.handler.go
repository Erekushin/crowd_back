package role

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type ActionHandler struct{}

func (*ActionHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = ActionList(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*ActionHandler) Create(c *fiber.Ctx) error {
	action := new(Action)

	if err := c.BodyParser(action); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*action); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	action.CreatedBy = oauth.GetSessionUserId(c)
	if err := action.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*ActionHandler) Update(c *fiber.Ctx) error {
	action := new(Action)

	if err := c.BodyParser(action); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*action); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	action.UpdatedBy = oauth.GetSessionUserId(c)
	if err := action.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*ActionHandler) Delete(c *fiber.Ctx) error {
	action := new(Action)

	if err := c.BodyParser(action); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*action); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	action.DeletedBy = oauth.GetSessionUserId(c)
	if err := action.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

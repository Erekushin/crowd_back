package role

import (
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"

	"crowdfund/pkg/core"
)

type RoleHandler struct{}

func (*RoleHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)
	orgId := convertor.StringToInt(c.Query("org_id"))

	if orgId == 0 {
		return core.Resolve(400, c, core.Response("org_id is required"))
	}
	if res, err = RoleList(c, orgId); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*RoleHandler) Create(c *fiber.Ctx) error {
	var role Role

	if err := c.BodyParser(&role); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(role); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	role.CreatedBy = oauth.GetSessionUserId(c)
	if err := role.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*RoleHandler) Update(c *fiber.Ctx) error {
	role := new(Role)

	if err := c.BodyParser(&role); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(role); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	role.UpdatedBy = oauth.GetSessionUserId(c)
	if err := role.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*RoleHandler) Delete(c *fiber.Ctx) error {
	role := new(Role)

	if err := c.BodyParser(role); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*role); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	role.DeletedBy = oauth.GetSessionUserId(c)
	if err := role.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

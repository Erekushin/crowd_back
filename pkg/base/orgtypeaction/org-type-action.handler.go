package orgtypeaction

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type OrgTypeActionHandler struct{}

func (*OrgTypeActionHandler) List(c *fiber.Ctx) error {
	orgId := convertor.StringToInt(c.Query("type_id"))
	if orgId == 0 {
		return core.Resolve(400, c, core.Response("type_id is required"))
	}
	res := OrgTypeActionList(c, orgId)
	return core.Resolve(200, c, core.Response("success", res))
}

func (*OrgTypeActionHandler) Add(c *fiber.Ctx) error {
	data := new(OrgTypeAction)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.CreatedBy = oauth.GetSessionUserId(c)
	if err := data.Add(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}
func (*OrgTypeActionHandler) Remove(c *fiber.Ctx) error {
	data := new(OrgTypeAction)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := data.Remove(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

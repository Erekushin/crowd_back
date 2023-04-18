package organization

import (
	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type OrgVehicleHandler struct{}

func (*OrgVehicleHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)
	orgId := convertor.StringToInt(c.Query("org_id"))
	if orgId == 0 {
		return core.Resolve(400, c, core.Response("org_id is required"))

	}

	if res, err = VehicleList(c, orgId); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*OrgVehicleHandler) Add(c *fiber.Ctx) error {
	data := new(OrgVehicle)

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
func (*OrgVehicleHandler) Remove(c *fiber.Ctx) error {
	data := new(OrgVehicle)

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

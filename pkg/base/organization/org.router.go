package organization

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var c OrganizationHandler
	v := app.Group("/organization", oauth.TokenMiddleware)
	v.Get("find", c.Find)
	v.Get("", c.List)
	v.Post("", c.Create)
	v.Put("", c.Update)
	v.Delete("", c.Delete)

	var ou OrgUserHandler
	v.Get("user", ou.List)
	v.Post("user", ou.Add)
	v.Delete("user", ou.Remove)

	var ot OrgTerminalHandler
	v.Get("terminal", ot.List)
	v.Post("terminal", ot.Add)
	v.Delete("terminal", ot.Remove)

	var ov OrgVehicleHandler
	v.Get("vehicle", ov.List)
	v.Post("vehicle", ov.Add)
	v.Delete("vehicle", ov.Remove)

	var ott OrgTypeHandler
	v.Get("type", ott.OrgTypeList)
}

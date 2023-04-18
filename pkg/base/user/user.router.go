package user

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var c UserHandler
	u := app.Group("/user", oauth.TokenMiddleware)
	u.Get("find", c.Find)
	u.Get("", c.List)
	u.Put("", c.Update)
	u.Delete("", c.Delete)
	u.Put("password", c.ChangePassword)
	u.Put("username", c.ChangeUsername)
	u.Get("find-phone", c.FindPhone)
	u.Post("identity/change", c.ChangeEmailOrPhone)
	var uo UserOrgHandler
	u.Get("org", uo.OrgList)
}

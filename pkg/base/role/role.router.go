package role

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var role RoleHandler
	r := app.Group("role", oauth.TokenMiddleware)
	r.Get("", role.List)
	r.Post("", role.Create)
	r.Put("", role.Update)
	r.Delete("", role.Delete)

	var module ModuleHandler
	m := app.Group("module", oauth.TokenMiddleware)
	m.Get("", module.List)
	m.Post("", module.Create)
	m.Put("", module.Update)
	m.Delete("", module.Delete)

	var page PageHandler
	p := app.Group("page", oauth.TokenMiddleware)
	p.Get("", page.List)
	p.Post("", page.Create)
	p.Put("", page.Update)
	p.Delete("", page.Delete)

	var action ActionHandler
	a := app.Group("action", oauth.TokenMiddleware)
	a.Get("", action.List)
	a.Post("", action.Create)
	a.Put("", action.Update)
	a.Delete("", action.Delete)
}

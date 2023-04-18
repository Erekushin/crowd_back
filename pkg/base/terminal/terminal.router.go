package terminal

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var c TerminalHandler
	v := app.Group("/terminal", oauth.TokenMiddleware)
	v.Get("find", c.Find)
	v.Post("", c.Create)
	v.Get("", c.List)
	v.Put("", c.Update)
	v.Delete("", c.Delete)
}

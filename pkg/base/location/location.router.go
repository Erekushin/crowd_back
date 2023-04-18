package location

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var location LocationHandler
	l := app.Group("location", oauth.TokenMiddleware)
	l.Get("", location.List)
	l.Post("", location.Create)
	l.Put("", location.Update)
	l.Delete("", location.Delete)
}

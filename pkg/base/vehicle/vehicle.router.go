package vehicle

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var c VehicleHandler
	v := app.Group("/vehicle", oauth.TokenMiddleware)
	v.Get("find", c.Find)
	v.Get("", c.List)
}

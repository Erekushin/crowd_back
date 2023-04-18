package country

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var country CountryHandler
	l := app.Group("countries", oauth.TokenMiddleware)
	l.Get("", country.List)
}

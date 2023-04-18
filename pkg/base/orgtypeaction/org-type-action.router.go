package orgtypeaction

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var c OrgTypeActionHandler
	v := app.Group("/org-type-action", oauth.TokenMiddleware)
	v.Get("", c.List)
	v.Post("", c.Add)
	v.Delete("", c.Remove)
}

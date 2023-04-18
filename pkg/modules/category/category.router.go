package category

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var handler CategoryHandler
	r := app.Group("category")
	r.Post("", handler.Create, oauth.TokenMiddleware)
	r.Put("", handler.Update, oauth.TokenMiddleware)
	r.Delete("", handler.Delete, oauth.TokenMiddleware)
	r.Get("", handler.List, oauth.TokenMiddleware)

}

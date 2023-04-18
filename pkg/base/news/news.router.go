package news

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var handler NewsHandler
	r := app.Group("news")
	r.Get("", handler.List)
	r.Post("", handler.Create, oauth.TokenMiddleware)
	r.Put("", handler.Update, oauth.TokenMiddleware)
	r.Delete("", handler.Delete, oauth.TokenMiddleware)

	var h NewsTypeHandler
	r.Get("type", h.List, oauth.TokenMiddleware)
	r.Post("type", h.Create, oauth.TokenMiddleware)
	r.Put("type", h.Update, oauth.TokenMiddleware)
	r.Delete("type", h.Delete, oauth.TokenMiddleware)
}

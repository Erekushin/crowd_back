package crowdfund

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var handler CrowdfundHandler
	r := app.Group("crowdfund")
	r.Get("", handler.List, oauth.TokenMiddleware)
	r.Post("", handler.Create, oauth.TokenMiddleware)
	r.Put("", handler.Update, oauth.TokenMiddleware)
	r.Delete("", handler.Delete, oauth.TokenMiddleware)
	r.Get("confirmed", handler.ListConfirmed)
	r.Post("confirm", handler.Confirm, oauth.TokenMiddleware)
	r.Post("cancel", handler.Cancel, oauth.TokenMiddleware)
	r.Post("info", handler.Info, oauth.TokenMiddleware)

}

package crowdfunduser

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var handler CrowdfundUserHandler
	r := app.Group("crowdfund_user", oauth.TokenMiddleware)
	r.Post("", handler.Create)
	r.Get("", handler.ListUser)
	r.Get("crowdfund", handler.ListUserCrowdfund)

}

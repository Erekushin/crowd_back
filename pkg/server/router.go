package server

import (
	"crowdfund/pkg/base/auth"
	"crowdfund/pkg/base/country"
	"crowdfund/pkg/base/document"
	"crowdfund/pkg/base/lang"
	"crowdfund/pkg/base/location"
	"crowdfund/pkg/base/news"
	"crowdfund/pkg/base/organization"
	"crowdfund/pkg/base/orgtypeaction"
	"crowdfund/pkg/base/role"
	"crowdfund/pkg/base/terminal"
	"crowdfund/pkg/base/user"
	"crowdfund/pkg/base/vehicle"
	"crowdfund/pkg/base/wallet"
	"crowdfund/pkg/modules/category"
	"crowdfund/pkg/modules/crowdfund"
	crowdfunduser "crowdfund/pkg/modules/crowdfund-user"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	// app.Use(module.RequestMiddleware)
	// core.SetRoutes(app)
	auth.SetRoutes(app)
	user.SetRoutes(app)
	role.SetRoutes(app)
	document.SetRoutes(app)
	organization.SetRoutes(app)
	terminal.SetRoutes(app)
	location.SetRoutes(app)
	vehicle.SetRoutes(app)
	country.SetRoutes(app)
	lang.SetRoutes(app)
	news.SetRoutes(app)
	orgtypeaction.SetRoutes(app)
	wallet.SetRoutes(app)
	crowdfund.SetRoutes(app)
	category.SetRoutes(app)
	crowdfunduser.SetRoutes(app)
}

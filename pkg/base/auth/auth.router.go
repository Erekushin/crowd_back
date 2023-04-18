package auth

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var a AuthHandler
	u := app.Group("auth")
	u.Get("identify", a.Identify)
	u.Post("login", a.Login)
	u.Post("register", a.Register)
	u.Put("org", oauth.TokenMiddleware, a.ChangeOrg)
	u.Get("password", a.ResetPasswordOtp)
	u.Put("password", a.ResetPassword)
	u.Get("qrlogin", oauth.TokenMiddleware, a.QrLogin)
}

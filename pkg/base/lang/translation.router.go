package lang

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var key KeyHandler
	k := app.Group("key", oauth.TokenMiddleware)
	k.Get("", key.List)
	k.Post("", key.Save)
	k.Put("", key.Update)
	k.Delete("", key.Delete)

	var lang LangHandler
	l := app.Group("language")
	l.Get("", lang.List)
	l.Post("", lang.Save, oauth.TokenMiddleware)
	l.Put("", lang.Update, oauth.TokenMiddleware)
	l.Delete("", lang.Delete, oauth.TokenMiddleware)

	var translation TranslationHandler
	t := app.Group("translation")
	t.Get("", translation.List)
	t.Put("", translation.Set, oauth.TokenMiddleware)
}

package document

import (
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	var docHandler DocumentHandler
	document := app.Group("document", oauth.TokenMiddleware)
	document.Post("/find", docHandler.Find)
	document.Post("/set_user", docHandler.SetUser)
	document.Get("", docHandler.List)
	document.Put("", docHandler.Update)
	document.Delete("", docHandler.Delete)

	var category CategoryHandler
	c := document.Group("category")
	c.Get("", category.List)
	c.Post("", category.Create, oauth.TokenMiddleware)
	c.Put("", category.Update, oauth.TokenMiddleware)
	c.Delete("", category.Delete, oauth.TokenMiddleware)

	var document_type TypeHandler
	t := document.Group("type")
	t.Get("", document_type.List)
	t.Post("", document_type.Create, oauth.TokenMiddleware)
	t.Put("", document_type.Update, oauth.TokenMiddleware)
	t.Delete("", document_type.Delete, oauth.TokenMiddleware)

}

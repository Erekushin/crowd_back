package lang

import (
	"crowdfund/pkg/core"

	"github.com/gofiber/fiber/v2"
)

type TranslationHandler struct{}

func (*TranslationHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = TranslationList(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*TranslationHandler) Set(c *fiber.Ctx) error {
	data := new(RequestTranslations)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	if err := data.TranslationSet(c); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

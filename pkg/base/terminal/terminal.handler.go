package terminal

import (
	"os"

	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type TerminalHandler struct{}

func (*TerminalHandler) Find(c *fiber.Ctx) error {
	searchText := c.Query("search_text")
	if searchText == "" {
		return core.Resolve(400, c, core.Response("search value is null"))
	}

	body := make(map[string]interface{})
	body["search_text"] = searchText

	resp := client.Request(os.Getenv("URL_FIND_TERMINAL"), "POST", body)

	return core.Resolve(200, c, core.Response("success", resp))
}

func (*TerminalHandler) Create(c *fiber.Ctx) error {
	data := new(Terminal)
	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.CreatedBy = oauth.GetSessionUserId(c)
	if err := data.Create(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*TerminalHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = List(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

func (*TerminalHandler) Update(c *fiber.Ctx) error {
	data := new(Terminal)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.UpdatedBy = oauth.GetSessionUserId(c)
	if err := data.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*TerminalHandler) Delete(c *fiber.Ctx) error {
	data := new(Terminal)

	if err := c.BodyParser(data); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*data); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	data.DeletedBy = oauth.GetSessionUserId(c)
	if err := data.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

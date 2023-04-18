package category

import (
	"fmt"

	"crowdfund/pkg/core"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
}

func (*CategoryHandler) Create(c *fiber.Ctx) error {
	category := new(Category)

	if err := c.BodyParser(category); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*category); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	category.CreatedBy = oauth.GetSessionUserId(c)
	if r, e := category.Create(); e != nil {
		category.Id = r.Id
		return core.Resolve(500, c, core.Response(e.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CategoryHandler) Update(c *fiber.Ctx) error {
	category := new(Category)

	if err := c.BodyParser(category); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*category); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}
	category.UpdatedBy = oauth.GetSessionUserId(c)
	if err := category.Update(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CategoryHandler) Delete(c *fiber.Ctx) error {
	category := new(Category)
	fmt.Println(".... ", c)
	if err := c.BodyParser(category); err != nil {
		return core.Resolve(400, c, core.Response(err.Error()))
	}

	if errors := core.Validate(*category); errors != nil {
		return core.Resolve(400, c, core.Response("validation error", errors))
	}

	category.DeletedBy = oauth.GetSessionUserId(c)
	if err := category.Delete(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response())
}

func (*CategoryHandler) List(c *fiber.Ctx) error {

	var (
		res interface{}
		err error
	)

	if res, err = List(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

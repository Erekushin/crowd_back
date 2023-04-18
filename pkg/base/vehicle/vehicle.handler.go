package vehicle

import (
	"os"

	"crowdfund/pkg/core"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type VehicleHandler struct{}

func (*VehicleHandler) Find(c *fiber.Ctx) error {
	searchText := c.Query("search_text")
	if searchText == "" {
		return core.Resolve(400, c, core.Response("search value is null"))
	}

	vehicle := Vehicle{}

	body := make(map[string]interface{})
	body["search_text"] = searchText

	resp := client.Request(os.Getenv("URL_FIND_VEHICLE"), "POST", body)
	if resp.Code != 200 {
		return c.JSON(resp)
	}

	convertor.MapToStruct(resp.Result, &vehicle)

	vehicle.CreatedBy = oauth.GetSessionUserId(c)
	if err := vehicle.Find(); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}

	return core.Resolve(200, c, core.Response("success", vehicle))
}

func (*VehicleHandler) List(c *fiber.Ctx) error {
	var (
		res interface{}
		err error
	)

	if res, err = List(c); err != nil {
		return core.Resolve(500, c, core.Response(err.Error()))
	}
	return core.Resolve(200, c, core.Response("success", res))
}

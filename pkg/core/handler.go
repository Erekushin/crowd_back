package core

import (
	"encoding/json"
	"os"

	"crowdfund/pkg/base/message"
	"crowdfund/pkg/helpers/client"
	"crowdfund/pkg/helpers/convertor"
	"crowdfund/pkg/oauth"

	"github.com/gofiber/fiber/v2"
)

type ServiceHandler struct{}

func SetRoutes(app *fiber.App) {
	var c ServiceHandler
	api := app.Group("/request", oauth.TokenMiddleware)
	api.Post("make", c.Make)
}

func (*ServiceHandler) Make(c *fiber.Ctx) error {

	msg := new(message.Message)
	msgCode := convertor.StringToInt(c.Get("message_code"))
	if msgCode == 0 {
		msg.Path = c.Path()
		msg.Method = c.Method()
		msg.FindByPath()
		if msg.Id != 0 {
			// go SaveRequestLog(c)
		}
	} else {
		msg.Code = uint(msgCode)
		msg.FindByCode()
	}

	if msg.Id == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "Service message not found"})
	}

	body := make(map[string]interface{})
	json.Unmarshal([]byte(c.Body()), &body)
	url := "http://localhost" + os.Getenv("PORT") + msg.Path
	headers := map[string]string{"Authorization": c.Get("Authorization")}
	resp := client.Request(url, msg.Method, body, headers)
	return c.JSON(resp)
}

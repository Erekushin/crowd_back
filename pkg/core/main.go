package core

import (
	"crowdfund/pkg/base/message"
	"crowdfund/pkg/helpers/convertor"

	"github.com/gofiber/fiber/v2"
)

func RequestMiddleware(c *fiber.Ctx) error {
	msg := new(message.Message)
	msgCode := convertor.StringToInt(c.Get("message_code"))
	if msgCode == 0 {
		msg.Path = c.Path()
		msg.Method = c.Method()
		msg.FindByPath()
		if msg.Id != 0 {
			go SaveRequestLog(c)
		}
	} else {
		msg.Code = uint(msgCode)
		msg.FindByCode()
	}

	if msg.Id == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "Service message not found"})
	}

	c.Locals("message", msg)
	return c.Next()
}

func Resolve(statusCode uint, c *fiber.Ctx, res *ApiResponse) error {
	// go SaveResponseLog(statusCode, c, res)
	return c.Status(int(statusCode)).JSON(res)
}

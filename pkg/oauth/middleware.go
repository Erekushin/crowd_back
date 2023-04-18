package oauth

import (
	"strings"
	"time"

	"crowdfund/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func TokenMiddleware(c *fiber.Ctx) error {

	header := string(c.Request().Header.Peek("Authorization"))

	if header == "" || !strings.Contains(header, "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON("invalid token")
	}
	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON("invalid token")
	}

	db := database.DBconn
	token := splitToken[1]

	session := &Session{}
	if err := db.Find(&session, "token = ? AND expires > ?", token, time.Now()).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}
	if session.Id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON("invalid or expired token")
	}

	c.Locals("session", session)
	return c.Next()
}

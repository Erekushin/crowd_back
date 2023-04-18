package server

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Init() {
	app := fiber.New()
	app.Use(requestid.New())
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New())

	InitRoutes(app)

	log.Fatal(app.Listen(os.Getenv("PORT")))
}

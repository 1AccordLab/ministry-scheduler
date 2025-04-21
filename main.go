package main

import (
	"log"
	"ministry-scheduler/views"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Static("/public", "./public")

	app.Static("/ui", "./spec")

	app.Get("/", render(views.Index("World!")))

	app.Get("/data", func(c *fiber.Ctx) error {
		return c.SendString("Hello, HTMX!")
	})

	log.Fatal(app.Listen(":3000"))
}

func render(component templ.Component) fiber.Handler {
	return adaptor.HTTPHandler(templ.Handler(component))
}

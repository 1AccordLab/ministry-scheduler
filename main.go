package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	app := fiber.New(
		fiber.Config{
			Views: html.New("./views", ".html"),
		})
	app.Static("/public", "./public")

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Name": "World!",
		})
	})

	app.Get("/data", func(c *fiber.Ctx) error {
		return c.SendString("Hello, HTMX!")
	})

	log.Fatal(app.Listen(":3000"))
}

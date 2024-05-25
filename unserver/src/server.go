package main

import "github.com/gofiber/fiber/v2"

func build_server() *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"error": false,
			"data":  true,
		})
	})

	app.Get("/transaction", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"error": false,
			"data": fiber.Map{
				"transactions": []string{},
			},
		})
	})

	return app
}

func main() {
	build_server().Listen(":3000")
}

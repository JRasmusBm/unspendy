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

	register_transaction_routes(app)

	return app
}

func main() {
	build_server().Listen(":3000")
}

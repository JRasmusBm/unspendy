package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func build_server(db *sql.DB) *fiber.App {
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

func start_server(db *sql.DB) {
	build_server(db).Listen(":3000")
}

func main() {
	with_db("temp.db", start_server)
}

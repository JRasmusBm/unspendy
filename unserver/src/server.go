package main

import "github.com/gofiber/fiber/v2"

type TransactionSearchResultData struct {
	Transactions []string `json:"transactions"`
}

type TransactionSearchResult struct {
	Error bool                        `json:"error"`
	Data  TransactionSearchResultData `json:"data"`
}

func build_server() *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"error": false,
			"data":  true,
		})
	})

	app.Get("/transaction", func(c *fiber.Ctx) error {
		return c.JSON(TransactionSearchResult{
			Error: false,
			Data: TransactionSearchResultData{
				Transactions: []string{},
			},
		})
	})

	return app
}

func main() {
	build_server().Listen(":3000")
}

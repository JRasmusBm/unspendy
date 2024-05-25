package main

import "github.com/gofiber/fiber/v2"

type TransactionSearchResultData struct {
	Transactions []string `json:"transactions"`
}

type TransactionSearchResult struct {
	Error bool                        `json:"error"`
	Data  TransactionSearchResultData `json:"data"`
}

func register_transaction_routes(app *fiber.App) {
	app.Get("/transaction", func(c *fiber.Ctx) error {
		return c.JSON(TransactionSearchResult{
			Error: false,
			Data: TransactionSearchResultData{
				Transactions: []string{},
			},
		})
	})
}

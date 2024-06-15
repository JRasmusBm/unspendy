package main

import "github.com/gofiber/fiber/v2"

type Transaction struct {
	transaction_date string
}

type TransactionSearchResultData struct {
	Transactions []Transaction `json:"transactions"`
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
				Transactions: []Transaction{},
			},
		})
	})
}

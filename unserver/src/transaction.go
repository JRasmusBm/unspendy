package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	transaction_date string
	// transaction_date: date = Field(validation_alias=AliasChoices("Date", "Datum"))
	// description: str = Field(
	//     validation_alias=AliasChoices("Name / Description", "Omschrijving")
	// )
	// account: str = Field(validation_alias=AliasChoices("Account", "Rekening"))
	// counterparty: str = Field(
	//     validation_alias=AliasChoices("Counterparty", "Rekening naam")
	// )
	// code: str = Field(validation_alias=AliasChoices("Code", "Tegenrekening"))
	// is_debit: bool = Field(validation_alias=AliasChoices("Debit/credit", "Af Bij"))
	// amount_in_cents: int = Field(
	//     validation_alias=AliasChoices("Amount (EUR)", "Bedrag")
	// )
	// transaction_type: str = Field(
	//     validation_alias=AliasChoices("Transaction type", "Mutatiesoort")
	// )
	// notifications: str = Field(
	//     validation_alias=AliasChoices("Notifications", "Mededelingen")
	// )
	// resulting_balance_in_cents: int = Field(
	//     validation_alias=AliasChoices("Resulting balance", "Saldo na mutatie")
	// )
}

type TransactionSearchResultData struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionSearchResult struct {
	Error bool                        `json:"error"`
	Data  TransactionSearchResultData `json:"data"`
}

func migrate_transactions(db *sql.DB) {
	db.Exec(`
	CREATE TABLE IF NOT EXISTS transactions (
    id UUID 
	);
		`)
}

func save_transaction(db *sql.DB) {

}

func respond_with_error(c *fiber.Ctx, err error) error {
	return c.Status(400).JSON(fiber.Map{
		"error":   true,
		"message": err,
	})
}

func register_transaction_routes(app *fiber.App, db *sql.DB) {
	migrate_transactions(db)

	app.Post("/transaction/upload", func(c *fiber.Ctx) error {
		file_header, err := c.FormFile("upload")
		if err != nil {
			return respond_with_error(c, err)
		}

		file, err := file_header.Open()
		defer file.Close()
		if err != nil {
			return respond_with_error(c, err)
		}

		reader := csv.NewReader(file)
		contents, err := reader.ReadAll()
		if err != nil {
			return respond_with_error(c, err)
		}

		fmt.Printf("%v\n", contents)

		return c.JSON(fiber.Map{
			"error": false,
		})
	})

	app.Get("/transaction", func(c *fiber.Ctx) error {
		return c.JSON(TransactionSearchResult{
			Error: false,
			Data: TransactionSearchResultData{
				Transactions: []Transaction{},
			},
		})
	})
}

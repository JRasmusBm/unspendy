package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	TransactionDate string `json:"transaction_date"`
	Description     string `json:"description"`
	Account         string `json:"account"`
	Counterparty    string `json:"counterparty"`
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

func migrate_transactions(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS transactions (
    id UUID NOT NULL PRIMARY KEY,
    transaction_date DATE,
    description TEXT,
    account TEXT,
    counterparty TEXT
	);`)
	return err
}

func save_transaction(db *sql.DB, t Transaction) error {
	_, err := db.Exec(`
INSERT INTO transactions (id, transaction_date, description, account, counterparty)
VALUES ($1, $2, $3, $4, $5);
	`, uuid.New(), t.TransactionDate, t.Description, t.Account, t.Counterparty)
	return err
}

func search_transactions(db *sql.DB) ([]Transaction, error) {
	rows, err := db.Query(`
	SELECT date(transaction_date), description, account, counterparty FROM transactions;
	`)
	defer rows.Close()

	result := make([]Transaction, 0)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var t Transaction
		err = rows.Scan(
			&t.TransactionDate,
			&t.Description,
			&t.Account,
			&t.Counterparty,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, err
}

type TransactionSearchResultData struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionSearchResult struct {
	Error bool                        `json:"error"`
	Data  TransactionSearchResultData `json:"data"`
}

type ErrorPayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func wrap_error(err error) ErrorPayload {
	return ErrorPayload{
		Error:   true,
		Message: fmt.Sprintf("error: %v", err),
	}
}

func register_transaction_routes(app *fiber.App, db *sql.DB) {
	err := migrate_transactions(db)

	if err != nil {
		log.Fatalf("%v\n", err)
		panic(err)
	}

	app.Post("/transaction/upload", func(c *fiber.Ctx) error {
		file_header, err := c.FormFile("upload")

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(wrap_error(err))
		}

		file, err := file_header.Open()
		defer file.Close()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(wrap_error(err))
		}

		reader := csv.NewReader(file)
		reader.Comma = ';'
		contents, err := reader.ReadAll()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(wrap_error(err))
		}

		header_row := contents[0]
		for key, row := range contents {
			if key == 0 {
				continue
			}

			t := Transaction{}
			for i, field_name := range header_row {
				if field_name == "Date" || field_name == "Datum" {
					t.TransactionDate = row[i]
				}

				if field_name == "Name / Description" || field_name == "Omschrijving" {
					t.Description = row[i]
				}

				if field_name == "Account" || field_name == "Rekening" {
					t.Account = row[i]
				}

				if field_name == "Counterparty" || field_name == "Rekening naam" {
					t.Counterparty = row[i]
				}
			}

			err := save_transaction(db, t)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(wrap_error(err))
			}
		}

		return c.JSON(fiber.Map{
			"error": false,
		})
	})

	app.Get("/transaction", func(c *fiber.Ctx) error {
		result, err := search_transactions(db)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(wrap_error(err))
		}

		return c.JSON(TransactionSearchResult{
			Error: false,
			Data: TransactionSearchResultData{
				Transactions: result,
			},
		})
	})
}

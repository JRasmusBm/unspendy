package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	Id  string `json:"id"`
	TransactionDate  string `json:"transaction_date"`
	Description      string `json:"description"`
	Account          string `json:"account"`
	Counterparty     string `json:"counterparty"`
	Code             string `json:"code"`
	IsDebit          bool   `json:"is_debit"`
	AmountInCents    int    `json:"amount_in_cents"`
	TransactionType  string `json:"transaction_type"`
	Notifications    string `json:"notifications"`
	ResultingBalance int    `json:"resulting_balance_in_cents"`
}

func migrate_transactions(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS transactions (
    id UUID NOT NULL PRIMARY KEY,
    transaction_date DATE,
    description TEXT,
    account TEXT,
    counterparty TEXT,
    code TEXT,
    is_debit BOOLEAN,
		amount_in_cents INTEGER,
    transaction_type TEXT,
    notifications TEXT,
		resulting_balance_in_cents INTEGER
	);`)
	return err
}

func save_transaction(db *sql.DB, t Transaction) error {
	_, err := db.Exec(`
INSERT INTO transactions (
	id,
	transaction_date,
	description,
	account,
	counterparty,
	code,
	is_debit,
	amount_in_cents,
	transaction_type,
	notifications,
	resulting_balance_in_cents
) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11
	);
	`,
		uuid.New(),
		t.TransactionDate,
		t.Description,
		t.Account,
		t.Counterparty,
		t.Code,
		t.IsDebit,
		t.AmountInCents,
		t.TransactionType,
		t.Notifications,
		t.ResultingBalance,
	)
	return err
}

func search_transactions(db *sql.DB) ([]Transaction, error) {
	rows, err := db.Query(`
	SELECT 
	id,
	date(transaction_date),
	description,
	account,
	counterparty,
	code,
	is_debit,
	amount_in_cents,
	transaction_type,
	notifications,
	resulting_balance_in_cents
	FROM transactions;
	`)
	defer rows.Close()

	result := make([]Transaction, 0)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var t Transaction
		err = rows.Scan(
			&t.Id,
			&t.TransactionDate,
			&t.Description,
			&t.Account,
			&t.Counterparty,
			&t.Code,
			&t.IsDebit,
			&t.AmountInCents,
			&t.TransactionType,
			&t.Notifications,
			&t.ResultingBalance,
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

func parse_money(value string) (int, error) {
	result, err := strconv.ParseInt(strings.ReplaceAll(value, ",", ""), 10, 64)
	return int(result), err
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

				if field_name == "Code" || field_name == "Tegenrekening" {
					t.Code = row[i]
				}

				if field_name == "Debit/credit" || field_name == "Af Bij" {
					t.IsDebit = row[i] == "Debit" || row[i] == "Bij"
				}

				if field_name == "Amount (EUR)" || field_name == "Bedrag" {
					value, err := parse_money(row[i])
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(wrap_error(err))
					}

					t.AmountInCents = value
				}

				if field_name == "Transaction type" || field_name == "Mutatiesoort" {
					t.TransactionType = row[i]
				}

				if field_name == "Notifications" || field_name == "Mededelingen" {
					t.Notifications = row[i]
				}

				if field_name == "Resulting balance" || field_name == "Saldo na mutatie" {
					value, err := parse_money(row[i])
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(wrap_error(err))
					}
					t.ResultingBalance = value
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

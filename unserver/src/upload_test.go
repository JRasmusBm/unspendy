package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func run_upload_csv(t *testing.T, server *fiber.App, csv_cells [][]string) {
	csv_rows := []string{}
	for _, row := range csv_cells {
		csv_rows = append(csv_rows, strings.Join(row, ";"))
	}

	csvReader := strings.NewReader(strings.Join(csv_rows, "\n"))
	upload_form_body := new(bytes.Buffer)
	mw := multipart.NewWriter(upload_form_body)
	w, err := mw.CreateFormFile("upload", "something.csv")
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))
	io.Copy(w, csvReader)
	mw.Close()

	upload_req := httptest.NewRequest("POST", "/transaction/upload", upload_form_body)
	upload_req.Header.Add("Content-Type", mw.FormDataContentType())

	resp, err := server.Test(upload_req, -1)
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	var upload_result ErrorPayload
	body, err := io.ReadAll(resp.Body)
	json.Unmarshal(body, &upload_result)

	assert.Equal(t, false, upload_result.Error, upload_result.Message)
	assert.Equal(t, 200, resp.StatusCode)
}

func run_search_transactions(t *testing.T, server *fiber.App) (search_result TransactionSearchResult) {
	req := httptest.NewRequest("GET", "/transaction", nil)
	resp, err := server.Test(req, -1)
	defer resp.Body.Close()

	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	json.Unmarshal(body, &search_result)

	return search_result
}

func NewTransactionSearchResult(transactions []Transaction) TransactionSearchResult {
	return TransactionSearchResult{
		Error: false,
		Data:  TransactionSearchResultData{Transactions: transactions},
	}
}

func test_server_setup(t *testing.T) (*fiber.App, *sql.DB) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	server := build_server(db)
	return server, db
}

func TestUpload(t *testing.T) {
	t.Run("Without upload returns empty data", func(t *testing.T) {
		server, db := test_server_setup(t)
		defer db.Close()

		search_result := run_search_transactions(t, server)

		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{}),
			search_result,
		)
	})

	t.Run("With upload returns data uploaded in English", func(t *testing.T) {
		server, db := test_server_setup(t)
		defer db.Close()

		run_upload_csv(t, server, [][]string{
			{
				"Date",
				"Name / Description",
				"Account",
				"Counterparty",
			},
			{
				"2024-01-12",
				"this is a description",
				"this is an account",
				"this is a counterparty",
			},
		})

		search_result := run_search_transactions(t, server)

		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{
				{
					TransactionDate: "2024-01-12",
					Description:     "this is a description",
					Account:         "this is an account",
					Counterparty:    "this is a counterparty",
				},
			}),
			search_result,
		)
	})

	t.Run("With upload returns uploaded data uploaded in Dutch", func(t *testing.T) {
		server, db := test_server_setup(t)
		defer db.Close()

		run_upload_csv(t, server, [][]string{
			{
				"Datum",
				"Omschrijving",
				"Rekening",
				"Rekening naam",
			},
			{
				"2024-01-12",
				"this is a description",
				"this is an account",
				"this is a counterparty",
			},
		})

		search_result := run_search_transactions(t, server)

		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{
				{
					TransactionDate: "2024-01-12",
					Description:     "this is a description",
					Account:         "this is an account",
					Counterparty:    "this is a counterparty",
				},
			}),
			search_result,
		)
	})
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func NewTransactionSearchResult(transactions []Transaction) TransactionSearchResult {
	return TransactionSearchResult{
		Error: false,
		Data:  TransactionSearchResultData{Transactions: transactions},
	}
}

func TestUpload(t *testing.T) {
	t.Run("Without upload returns empty data", func(t *testing.T) {
		db, err := sql.Open("sqlite3", "./temp_test.db")
		defer db.Close()
		req := httptest.NewRequest("GET", "/transaction", nil)
		resp, _ := build_server(db).Test(req, -1)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		var transaction_search_result TransactionSearchResult
		json.Unmarshal(body, &transaction_search_result)
		fmt.Printf("%#v\n", transaction_search_result)

		assert.Equal(t, nil, err)
		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{}),
			transaction_search_result,
		)
	})

	t.Run("With upload returns uploaded data", func(t *testing.T) {
		db, err := sql.Open("sqlite3", "./temp_test.db")
		defer db.Close()
		server := build_server(db)
		upload_resp, _ := server.Test(httptest.NewRequest("POST", "/transaction/upload", nil), -1)
		assert.Equal(t, 200, upload_resp.StatusCode)
		resp, _ := server.Test(httptest.NewRequest("GET", "/transaction", nil), -1)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		var transaction_search_result TransactionSearchResult
		json.Unmarshal(body, &transaction_search_result)
		fmt.Printf("%#v\n", transaction_search_result)

		assert.Equal(t, nil, err)
		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{
				{transaction_date: "2024-05-13T09:14:22.000Z"},
			}),
			transaction_search_result,
		)
	})
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func get_result[T any](response *http.Response) (int, T) {
	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	var result T
	json.Unmarshal(body, &result)
	return response.StatusCode, result
}

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

		csvReader := strings.NewReader(`
		name,age,city
		Alice,30,New York
		Bob,25,Los Angeles
		`)

		server := build_server(db)
		upload_resp, _ := server.Test(httptest.NewRequest("POST", "/transaction/upload", csvReader), -1)
		upload_status, upload_body := get_result[map[string]bool](upload_resp)
		fmt.Printf("upload_body: %v\n", upload_body)
		assert.Equal(t, true, upload_body["error"] )
		assert.Equal(t, 200, upload_status)

		resp, _ := server.Test(httptest.NewRequest("GET", "/transaction", nil), -1)
	  status, body := get_result[TransactionSearchResult](resp)
		fmt.Printf("%#v\n", body)

		assert.Equal(t, nil, err)
		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{
				{transaction_date: "2024-05-13T09:14:22.000Z"},
			}),
			body,
		)
		assert.Equal(t, 200, status)
	})
}

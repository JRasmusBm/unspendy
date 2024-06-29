package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
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
Date
2024-01-12
`)

		upload_form_body := new(bytes.Buffer)
		mw := multipart.NewWriter(upload_form_body)
		w, err := mw.CreateFormFile("upload", "something.csv")
		io.Copy(w, csvReader)
		mw.Close()
		server := build_server(db)
		upload_req := httptest.NewRequest("POST", "/transaction/upload", upload_form_body)
		upload_req.Header.Add("Content-Type", mw.FormDataContentType())
		upload_resp, _ := server.Test(upload_req, -1)
		upload_status, upload_body := get_result[ErrorPayload](upload_resp)

		assert.Equal(t, false, upload_body.Error, upload_body.Message)
		assert.Equal(t, 200, upload_status)

		resp, _ := server.Test(httptest.NewRequest("GET", "/transaction", nil), -1)
		status, body := get_result[TransactionSearchResult](resp)

		assert.Equal(t, nil, err, err)
		assert.Equal(t,
			NewTransactionSearchResult([]Transaction{
				{transaction_date: "2024-01-12"},
			}),
			body,
		)
		assert.Equal(t, 200, status)
	})
}

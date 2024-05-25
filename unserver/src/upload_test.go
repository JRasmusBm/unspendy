package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UploadResult struct {
	Error bool `json:"error"`
	Data  bool `json:"data"`
}

type TransactionSearchResultData struct {
	Transactions []string `json:"transactions"`
}

type TransactionSearchResult struct {
	Error bool                        `json:"error"`
	Data  TransactionSearchResultData `json:"data"`
}

func NewTransactionSearchResult(transactions []string) TransactionSearchResult {
	return TransactionSearchResult{
		Error: false,
		Data:  TransactionSearchResultData{Transactions: transactions},
	}
}

func TestUpload(t *testing.T) {
	t.Run("Without upload returns empty data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/transaction", nil)
		resp, _ := build_server().Test(req, -1)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		var transaction_search_result TransactionSearchResult
		json.Unmarshal(body, &transaction_search_result)
		fmt.Printf("%#v\n", transaction_search_result)

		assert.Equal(t, nil, err)
		assert.Equal(t,
			NewTransactionSearchResult([]string{}),
			transaction_search_result,
		)
	})

}

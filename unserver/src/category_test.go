package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func NewCategorySearchResult(categories []Category) CategorySearchResult {
	return CategorySearchResult{
		Error: false,
		Data:  CategorySearchResultData{Categories: categories},
	}
}

func run_create_category(t *testing.T, server *fiber.App, c Category) error {
	encoded, err := json.Marshal(c)
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	req := httptest.NewRequest("POST", "/category", bytes.NewBuffer(encoded))
	resp, err := server.Test(req, -1)
	defer resp.Body.Close()

	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	var payload ErrorPayload
	json.Unmarshal(body, &payload)

	assert.Equal(t, false, payload.Error, payload.Message)
	assert.Equal(t, 200, resp.StatusCode)

	return err
}

func run_search_categories(t *testing.T, server *fiber.App) (search_result CategorySearchResult) {
	req := httptest.NewRequest("GET", "/category", nil)
	resp, err := server.Test(req, -1)
	defer resp.Body.Close()

	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	err = json.Unmarshal(body, &search_result)
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	return search_result
}

func TestCategory(t *testing.T) {
	t.Run("Without upload returns empty data", func(t *testing.T) {
		server, db := test_server_setup(t)
		defer db.Close()

		search_result := run_search_categories(t, server)

		assert.Equal(t,
			NewCategorySearchResult([]Category{}),
			search_result,
		)
	})

	t.Run("Stores created category", func(t *testing.T) {
		server, db := test_server_setup(t)
		defer db.Close()

		run_create_category(t, server, Category{
			Name: "The name of the category",
		})
		search_result := run_search_categories(t, server)

		assert.NotEmpty(t, search_result.Data.Categories[0].Id)
		search_result.Data.Categories[0].Id = "override"

		assert.Equal(t,
			NewCategorySearchResult([]Category{{
				Id:   "override",
				Name: "The name of the category",
			}}),
			search_result,
		)
	})
}

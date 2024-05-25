package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseJSON(r io.ReadCloser, v interface{}) error {
	// Read the entire content into a byte slice
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the provided interface
	return json.Unmarshal(data, v)
}

type HealthResult struct {
	Error bool `json:"error"`
	Data  bool `json:"data"`
}

func TestHealthCheck(t *testing.T) {
	t.Run("Returns health as true", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, _ := build_server().Test(req, -1)
		defer resp.Body.Close()

		var health_result HealthResult
		err := parseJSON(resp.Body, &health_result)

		assert.Equal(t, nil, err)
		assert.Equal(t, false, health_result.Error)
		assert.Equal(t, true, health_result.Data)
	})

}

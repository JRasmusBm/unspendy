package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HealthResult struct {
	Error bool `json:"error"`
	Data  bool `json:"data"`
}

func TestHealthCheck(t *testing.T) {
	t.Run("Returns health as true", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, _ := build_server().Test(req, -1)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		var health_result HealthResult
		json.Unmarshal(body, &health_result)

		assert.Equal(t, nil, err)
		assert.Equal(t, false, health_result.Error)
		assert.Equal(t, true, health_result.Data)
	})

}

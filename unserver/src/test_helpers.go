package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func test_server_setup(t *testing.T) (*fiber.App, *sql.DB) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.Equal(t, nil, err, fmt.Sprintf("%#v", err))

	server := build_server(db)
	return server, db
}

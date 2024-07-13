package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type Category struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func migrate_categories(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS categories (
    id UUID NOT NULL PRIMARY KEY,
    name TEXT
	);`)
	return err
}

func save_category(db *sql.DB, c Category) error {
	_, err := db.Exec(`
INSERT INTO categories (
	id,
	name
) VALUES (
		$1,
		$2
	);
	`,
		uuid.New(),
		c.Name,
	)
	return err
}

func search_categories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query(`
	SELECT 
	id,
	name
	FROM categories;
	`)
	defer rows.Close()

	result := make([]Category, 0)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var c Category
		err = rows.Scan(
			&c.Id,
			&c.Name,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, err
}

type CategorySearchResultData struct {
	Categories []Category `json:"categories"`
}

type CategorySearchResult struct {
	Error bool                     `json:"error"`
	Data  CategorySearchResultData `json:"data"`
}

func register_category_routes(app *fiber.App, db *sql.DB) {
	err := migrate_categories(db)

	if err != nil {
		log.Fatalf("%v\n", err)
		panic(err)
	}

	app.Post("/category", func(c *fiber.Ctx) error {
		var category Category
		err := json.Unmarshal(c.Body(), &category)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(wrap_error(err))
		}

		save_category(db, category)

		return c.JSON(fiber.Map{
			"error": false,
		})
	})

	app.Get("/category", func(c *fiber.Ctx) error {
		result, err := search_categories(db)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(wrap_error(err))
		}

		return c.JSON(CategorySearchResult{
			Error: false,
			Data: CategorySearchResultData{
				Categories: result,
			},
		})
	})
}

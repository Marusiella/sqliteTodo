package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	app := fiber.New()
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	app.Use(logger.New())
	SetupRoutes(db, app)
	err = SetupDatabase(db)
	if err != nil {
		return
	}
	log.Fatal(app.Listen("0.0.0.0:3000"))
}

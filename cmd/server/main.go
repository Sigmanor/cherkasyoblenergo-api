package main

import (
	"log"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/database"
	"cherkasyoblenergo-api/internal/handlers"
	"cherkasyoblenergo-api/internal/middleware"
	"cherkasyoblenergo-api/internal/models"
	"cherkasyoblenergo-api/internal/parser"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Println("Failed to load config:", err)
		return
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("Failed to connect to DB:", err)
		return
	}

	db.AutoMigrate(&models.Schedule{}, &models.APIKey{})

	parser.StartCron(db)

	app := fiber.New()
	app.Use(middleware.APIKeyAuth(db))
	app.Use(middleware.RateLimiter(db))
	app.Use(middleware.Logger())

	api := app.Group("/cherkasyoblenergo/api")

	api.Post("/blackout-schedule", handlers.PostSchedule(db))
	api.Get("/generate-api-key", handlers.GenerateAPIKey(db, cfg))
	api.Get("/update-api-key", handlers.ManageAPIKey(db, cfg))

	log.Printf("Server is running on port %s (app version: %s)",
		cfg.SERVER_PORT,
		cfg.APP_VERSION)

	app.Listen(":" + cfg.SERVER_PORT)
}

package main

import (
	"fmt"
	"log"

	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/database"
	"cherkasyoblenergo-api/internal/handlers"
	"cherkasyoblenergo-api/internal/middleware"
	"cherkasyoblenergo-api/internal/models"
	"cherkasyoblenergo-api/internal/parser"

	"github.com/gofiber/fiber/v2"
)

func runServer() error {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return err
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database, got error %w", err)
	}

	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&models.Schedule{}, &models.APIKey{}); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}
	log.Println("Database migrations completed successfully")
	parser.StartCron(db)

	app := fiber.New()
	app.Use(middleware.APIKeyAuth(db))
	app.Use(middleware.RateLimiter(db))
	app.Use(middleware.Logger())

	api := app.Group("/cherkasyoblenergo/api")
	api.Post("/blackout-schedule", handlers.PostSchedule(db))
	api.Get("/generate-api-key", handlers.GenerateAPIKey(db, cfg))
	api.Get("/update-api-key", handlers.ManageAPIKey(db, cfg))

	log.Printf("Starting server (app version: %s)\n", config.AppVersion)
	return app.Listen(":" + cfg.ServerPort)
}

func main() {
	if err := runServer(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

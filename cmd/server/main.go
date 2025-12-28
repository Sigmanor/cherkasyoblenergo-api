package main

import (
	"cherkasyoblenergo-api/internal/cache"
	"cherkasyoblenergo-api/internal/config"
	"cherkasyoblenergo-api/internal/database"
	"cherkasyoblenergo-api/internal/handlers"
	"cherkasyoblenergo-api/internal/logger"
	"cherkasyoblenergo-api/internal/middleware"
	"cherkasyoblenergo-api/internal/models"
	"cherkasyoblenergo-api/internal/parser"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func runServer() error {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return err
	}

	logger.SetupGlobal(cfg.LogLevel)

	db, err := database.ConnectDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database, got error %w", err)
	}

	log.Println("Running database migrations...")
	if err := db.Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Silent)}).AutoMigrate(&models.Schedule{}, &models.IPRateLimit{}); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}
	log.Println("Database migrations completed successfully")

	cronScheduler := parser.StartCron(db, cfg.NewsURL)

	scheduleCache := cache.NewScheduleCache(cfg.CacheTTLSeconds)

	rateLimiter := middleware.NewIPRateLimiter(db, cfg.RateLimitPerMinute)

	fiberConfig := fiber.Config{}
	switch cfg.ProxyMode {
	case "cloudflare":
		fiberConfig.EnableTrustedProxyCheck = true
		fiberConfig.TrustedProxies = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1"}
		fiberConfig.ProxyHeader = "CF-Connecting-IP"
		log.Println("Proxy mode: cloudflare (using CF-Connecting-IP header)")
	case "standard":
		fiberConfig.EnableTrustedProxyCheck = true
		fiberConfig.TrustedProxies = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1"}
		fiberConfig.ProxyHeader = fiber.HeaderXForwardedFor
		log.Println("Proxy mode: standard (using X-Forwarded-For header)")
	default:
		log.Println("Proxy mode: none (using direct connection IP)")
	}

	app := fiber.New(fiberConfig)
	app.Use(fiberrecover.New())
	app.Use(middleware.HTTPSEnforcement())
	app.Use(middleware.APIKeyAuth(cfg.APIKey))
	app.Use(rateLimiter.Middleware())
	app.Use(middleware.Logger())

	api := app.Group("/cherkasyoblenergo/api")
	api.Get("/", handlers.GetAPIInfo())
	api.Get("/blackout-schedule", handlers.GetSchedule(db, scheduleCache))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		log.Println("Server started, running initial news parsing")
		go func() {
			parseCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			defer cancel()
			parser.FetchAndStoreNews(parseCtx, db, cfg.NewsURL)
		}()
		return nil
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()

		rateLimiter.Stop()
		log.Println("Rate limiter stopped")

		if cronScheduler != nil {
			cronScheduler.Stop()
			log.Println("Cron scheduler stopped")
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Printf("Starting server (app version: %s)\n", config.AppVersion)
	return app.Listen(":" + cfg.ServerPort)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := runServer(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

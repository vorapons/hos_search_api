package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"pt_search_hos/config"
	"pt_search_hos/handler"
	"pt_search_hos/repository"
	"pt_search_hos/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env (ignore error if file doesn't exist in prod)
	_ = godotenv.Load()

	cfg := config.Load()

	// Database
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Wire up layers
	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	staffSvc := services.NewStaffService(staffRepo, cfg.JWTSecret)
	patientSvc := services.NewPatientService(patientRepo)

	// Load persisted token blacklist from DB
	if err := staffSvc.LoadBlacklist(); err != nil {
		log.Printf("warning: could not load token blacklist: %v", err)
	}

	staffH := handler.NewStaffHandler(staffSvc)
	patientH := handler.NewPatientHandler(patientSvc)

	// Fiber
	app := fiber.New()
	app.Use(logger.New())
	handler.SetupRoutes(app, staffH, patientH, cfg.JWTSecret, staffSvc.IsTokenBlacklisted)

	// Start server in background
	go func() {
		log.Printf("server starting on :%s", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}

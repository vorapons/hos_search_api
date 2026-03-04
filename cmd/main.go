package main

import (
	"log"

	"pt_search_hos/config"
	"pt_search_hos/handler"
	"pt_search_hos/repository"
	"pt_search_hos/services"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
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

	staffH := handler.NewStaffHandler(staffSvc)
	patientH := handler.NewPatientHandler(patientSvc)

	// Fiber
	app := fiber.New()
	handler.SetupRoutes(app, staffH, patientH, cfg.JWTSecret, staffSvc.IsTokenBlacklisted)

	log.Printf("server starting on :%s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}

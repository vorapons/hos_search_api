package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pt_search_hos/config"
	ginhandler "pt_search_hos/handler/gin"
	"pt_search_hos/repository"
	"pt_search_hos/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	staffRepo   := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	staffSvc   := services.NewStaffService(staffRepo, cfg.JWTSecret)
	patientSvc := services.NewPatientService(patientRepo)

	// Load persisted token blacklist from DB
	if err := staffSvc.LoadBlacklist(); err != nil {
		log.Printf("warning: could not load token blacklist: %v", err)
	}

	staffH   := ginhandler.NewStaffHandler(staffSvc)
	patientH := ginhandler.NewPatientHandler(patientSvc)

	// Gin
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://hossearchapi-1071770156665.asia-southeast3.run.app", "http://localhost:3456", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.Use(gin.Logger(), gin.Recovery())
	ginhandler.SetupRoutes(r, staffH, patientH, cfg.JWTSecret, staffSvc.IsTokenBlacklisted)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	// Start server in background
	go func() {
		log.Printf("gin server starting on :%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gin server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("gin server stopped")
}

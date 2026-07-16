package main

import (
	"flag"
	"log"
	"net/http"

	"nalakarsa/internal/config"
	"nalakarsa/internal/database"
	"nalakarsa/internal/handler"
	"nalakarsa/internal/repository"
	"nalakarsa/internal/routes"
	"nalakarsa/internal/service"
	"nalakarsa/seed"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse CLI flags
	seedFlag := flag.Bool("seed", false, "run database seeder and exit")
	flag.Parse()

	log.Println("Starting Nalakarsa Backend API Server...")

	// 1. Load Configurations
	cfg := config.LoadConfig()

	// 2. Initialize Database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Critical Error: Database initialization failed: %v", err)
	}

	// 3. Check for Seed Flag
	if *seedFlag {
		if err := seed.SeedData(db); err != nil {
			log.Fatalf("Critical Error: Seeding database failed: %v", err)
		}
		log.Println("Database seeding completed. Exiting application...")
		return
	}

	// 4. Initialize Dependency Layers (Dependency Injection)
	// Repositories
	userRepo := repository.NewUserRepository(db)
	discRepo := repository.NewDiscussionRepository(db)
	collabRepo := repository.NewCollaborationRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	discService := service.NewDiscussionService(discRepo, userRepo)
	collabService := service.NewCollaborationService(collabRepo, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, cfg)
	discHandler := handler.NewDiscussionHandler(discService)
	collabHandler := handler.NewCollaborationHandler(collabService)

	// 5. Setup Routes
	r := routes.SetupRouter(cfg, authHandler, userHandler, discHandler, collabHandler)

	// Root Healthcheck endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "Nalakarsa Backend API Server",
			"version": "1.0.0",
		})
	})

	// 6. Start Server
	log.Printf("Server successfully started on port %s in %s mode\n", cfg.Port, cfg.Env)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Critical Error: Server failed to start: %v", err)
	}
}

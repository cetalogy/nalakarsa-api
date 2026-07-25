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
	connRepo := repository.NewConnectionRepository(db)
	projRepo := repository.NewProjectRepository(db)
	convRepo := repository.NewConversationRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, connRepo, projRepo)
	discService := service.NewDiscussionService(discRepo, userRepo)
	connService := service.NewConnectionService(connRepo, userRepo)
	projService := service.NewProjectService(projRepo, userRepo)
	convService := service.NewConversationService(convRepo, userRepo)
	notifService := service.NewNotificationService(notifRepo)
	dashService := service.NewDashboardService(projRepo, connRepo, convRepo, notifRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, cfg)
	discHandler := handler.NewDiscussionHandler(discService)
	connHandler := handler.NewConnectionHandler(connService)
	projHandler := handler.NewProjectHandler(projService)
	convHandler := handler.NewConversationHandler(convService)
	notifHandler := handler.NewNotificationHandler(notifService)
	dashHandler := handler.NewDashboardHandler(dashService)

	// 5. Setup Routes
	r := routes.SetupRouter(
		cfg,
		authHandler, userHandler, discHandler,
		projHandler, connHandler, convHandler,
		notifHandler, dashHandler,
	)

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

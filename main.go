package main

import (
	"flag"
	"log"
	"net/http"

	"nalakarsa/internal/config"
	"nalakarsa/internal/database"
	authhandler "nalakarsa/internal/handler/auth"
	connectionhandler "nalakarsa/internal/handler/connection"
	conversationhandler "nalakarsa/internal/handler/conversation"
	dashboardhandler "nalakarsa/internal/handler/dashboard"
	discussionhandler "nalakarsa/internal/handler/discussion"
	homepagehandler "nalakarsa/internal/handler/homepage"
	notificationhandler "nalakarsa/internal/handler/notification"
	projecthandler "nalakarsa/internal/handler/project"
	userhandler "nalakarsa/internal/handler/user"
	"nalakarsa/internal/middleware"
	connectionrepository "nalakarsa/internal/repository/connection"
	conversationrepository "nalakarsa/internal/repository/conversation"
	discussionrepository "nalakarsa/internal/repository/discussion"
	homerepository "nalakarsa/internal/repository/homepage"
	notificationrepository "nalakarsa/internal/repository/notification"
	projectrepository "nalakarsa/internal/repository/project"
	userrepository "nalakarsa/internal/repository/user"
	authservice "nalakarsa/internal/service/auth"
	connectionservice "nalakarsa/internal/service/connection"
	conversationservice "nalakarsa/internal/service/conversation"
	dashboardservice "nalakarsa/internal/service/dashboard"
	discussionservice "nalakarsa/internal/service/discussion"
	homepageService "nalakarsa/internal/service/homepage"
	notificationservice "nalakarsa/internal/service/notification"
	projectservice "nalakarsa/internal/service/project"
	userservice "nalakarsa/internal/service/user"

	"nalakarsa/internal/routes"
	authroutes "nalakarsa/internal/routes/auth"
	connectionroutes "nalakarsa/internal/routes/connection"
	conversationroutes "nalakarsa/internal/routes/conversation"
	dashboardroutes "nalakarsa/internal/routes/dashboard"
	discussionroutes "nalakarsa/internal/routes/discussion"
	homeroutes "nalakarsa/internal/routes/homepage"
	notificationroutes "nalakarsa/internal/routes/notification"
	projectroutes "nalakarsa/internal/routes/project"
	userroutes "nalakarsa/internal/routes/user"
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
	userRepo := userrepository.NewUserRepository(db)
	discRepo := discussionrepository.NewDiscussionRepository(db)
	connRepo := connectionrepository.NewConnectionRepository(db)
	projRepo := projectrepository.NewProjectRepository(db)
	convRepo := conversationrepository.NewConversationRepository(db)
	notifRepo := notificationrepository.NewNotificationRepository(db)
	homeRepo := homerepository.NewHomepageRepository(db)

	// Services
	authService := authservice.NewAuthService(userRepo, cfg)
	userService := userservice.NewUserService(userRepo, connRepo, projRepo)
	discService := discussionservice.NewDiscussionService(discRepo, userRepo)
	connService := connectionservice.NewConnectionService(connRepo, userRepo)
	projService := projectservice.NewProjectService(projRepo, userRepo, discRepo)
	convService := conversationservice.NewConversationService(convRepo, userRepo)
	notifService := notificationservice.NewNotificationService(notifRepo)
	dashService := dashboardservice.NewDashboardService(projRepo, connRepo, convRepo, notifRepo)
	homeService := homepageService.NewHomepageService(homeRepo)

	// Handlers
	authHandler := authhandler.NewAuthHandler(authService)
	userHandler := userhandler.NewUserHandler(userService, cfg)
	discHandler := discussionhandler.NewDiscussionHandler(discService)
	connHandler := connectionhandler.NewConnectionHandler(connService)
	projHandler := projecthandler.NewProjectHandler(projService)
	convHandler := conversationhandler.NewConversationHandler(convService)
	notifHandler := notificationhandler.NewNotificationHandler(notifService)
	dashHandler := dashboardhandler.NewDashboardHandler(dashService)
	homeHandler := homepagehandler.NewHomepageHandler(homeService)

	// 5. Setup Routes
	r := routes.NewRouter(cfg)
	v1 := r.Group("/api/v1")
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))

	authroutes.RegisterRoutes(v1, protected, authHandler)
	userroutes.RegisterRoutes(v1, protected, userHandler, connHandler.GetSuggestions)
	discussionroutes.RegisterRoutes(v1, protected, discHandler)
	connectionroutes.RegisterRoutes(protected, connHandler)
	projectroutes.RegisterRoutes(v1, protected, projHandler)
	conversationroutes.RegisterRoutes(protected, convHandler)
	notificationroutes.RegisterRoutes(protected, notifHandler)
	dashboardroutes.RegisterRoutes(protected, dashHandler)
	homeroutes.RegisterRoutes(v1, homeHandler)

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

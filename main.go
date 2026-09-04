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
	expertisehandler "nalakarsa/internal/handler/expertise"
	homepagehandler "nalakarsa/internal/handler/homepage"
	industryhandler "nalakarsa/internal/handler/industry"
	institutionhandler "nalakarsa/internal/handler/institution"
	knowledgehandler "nalakarsa/internal/handler/knowledge"
	locationhandler "nalakarsa/internal/handler/location"
	notificationhandler "nalakarsa/internal/handler/notification"
	organizationhandler "nalakarsa/internal/handler/organization"
	projecthandler "nalakarsa/internal/handler/project"
	userhandler "nalakarsa/internal/handler/user"
	"nalakarsa/internal/middleware"
	connectionrepository "nalakarsa/internal/repository/connection"
	conversationrepository "nalakarsa/internal/repository/conversation"
	discussionrepository "nalakarsa/internal/repository/discussion"
	expertiserepository "nalakarsa/internal/repository/expertise"
	homerepository "nalakarsa/internal/repository/homepage"
	industryrepository "nalakarsa/internal/repository/industry"
	institutionrepository "nalakarsa/internal/repository/institution"
	knowledgerepository "nalakarsa/internal/repository/knowledge"
	locationrepository "nalakarsa/internal/repository/location"
	notificationrepository "nalakarsa/internal/repository/notification"
	organizationrepository "nalakarsa/internal/repository/organization"
	projectrepository "nalakarsa/internal/repository/project"
	userrepository "nalakarsa/internal/repository/user"
	authservice "nalakarsa/internal/service/auth"
	connectionservice "nalakarsa/internal/service/connection"
	conversationservice "nalakarsa/internal/service/conversation"
	dashboardservice "nalakarsa/internal/service/dashboard"
	discussionservice "nalakarsa/internal/service/discussion"
	expertiseservice "nalakarsa/internal/service/expertise"
	homepageService "nalakarsa/internal/service/homepage"
	industryservice "nalakarsa/internal/service/industry"
	institutionservice "nalakarsa/internal/service/institution"
	knowledgeservice "nalakarsa/internal/service/knowledge"
	locationservice "nalakarsa/internal/service/location"
	notificationservice "nalakarsa/internal/service/notification"
	organizationservice "nalakarsa/internal/service/organization"
	projectservice "nalakarsa/internal/service/project"
	userservice "nalakarsa/internal/service/user"
	"nalakarsa/internal/utils"

	"nalakarsa/internal/routes"
	authroutes "nalakarsa/internal/routes/auth"
	connectionroutes "nalakarsa/internal/routes/connection"
	conversationroutes "nalakarsa/internal/routes/conversation"
	dashboardroutes "nalakarsa/internal/routes/dashboard"
	discussionroutes "nalakarsa/internal/routes/discussion"
	expertiseroutes "nalakarsa/internal/routes/expertise"
	homeroutes "nalakarsa/internal/routes/homepage"
	industryroutes "nalakarsa/internal/routes/industry"
	institutionroutes "nalakarsa/internal/routes/institution"
	knowledgeroutes "nalakarsa/internal/routes/knowledge"
	locationroutes "nalakarsa/internal/routes/location"
	notificationroutes "nalakarsa/internal/routes/notification"
	organizationroutes "nalakarsa/internal/routes/organization"
	projectroutes "nalakarsa/internal/routes/project"
	userroutes "nalakarsa/internal/routes/user"
	"nalakarsa/seed"

	"github.com/gin-gonic/gin"
)

func main() {
	seedFlag := flag.Bool("seed", false, "run database seeder and exit")
	flag.Parse()

	log.Println("Starting Nalakarsa Backend API Server...")
	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Critical Error: Database initialization failed: %v", err)
	}
	if *seedFlag {
		if err := seed.SeedData(db); err != nil {
			log.Fatalf("Critical Error: Seeding database failed: %v", err)
		}
		log.Println("Database seeding completed. Exiting application...")
		return
	}
	database.StartRefreshTokenCleanup(db)
	userRepo := userrepository.NewUserRepository(db)
	discRepo := discussionrepository.NewDiscussionRepository(db)
	expertiseRepo := expertiserepository.NewExpertiseRepository(db)
	connRepo := connectionrepository.NewConnectionRepository(db)
	projRepo := projectrepository.NewProjectRepository(db)
	convRepo := conversationrepository.NewConversationRepository(db)
	notifRepo := notificationrepository.NewNotificationRepository(db)
	homeRepo := homerepository.NewHomepageRepository(db)
	institutionRepo := institutionrepository.NewInstitutionRepository(db)
	organizationRepo := organizationrepository.NewOrganizationRepository(db)
	locationRepo := locationrepository.NewLocationRepository()
	industryRepo := industryrepository.NewIndustryRepository(db)
	knowledgeRepo := knowledgerepository.NewKnowledgeRepository(db)
	authService := authservice.NewAuthService(userRepo, cfg)
	userService := userservice.NewUserService(userRepo, connRepo, projRepo, notifRepo)
	discService := discussionservice.NewDiscussionService(discRepo, userRepo)
	expertiseService := expertiseservice.NewExpertiseService(expertiseRepo)
	connService := connectionservice.NewConnectionService(connRepo, userRepo)
	projService := projectservice.NewProjectService(projRepo, userRepo, discRepo, notifRepo)
	convService := conversationservice.NewConversationService(convRepo, userRepo, cfg)
	notifService := notificationservice.NewNotificationService(notifRepo)
	dashService := dashboardservice.NewDashboardService(projRepo, connRepo, convRepo, notifRepo)
	homeService := homepageService.NewHomepageService(homeRepo)
	institutionService := institutionservice.NewInstitutionService(institutionRepo)
	organizationService := organizationservice.NewOrganizationService(organizationRepo)
	locationService := locationservice.NewLocationService(locationRepo)
	industryService := industryservice.NewIndustryService(industryRepo)
	knowledgeService := knowledgeservice.NewKnowledgeService(knowledgeRepo)
	authHandler := authhandler.NewAuthHandler(authService)
	userHandler := userhandler.NewUserHandler(userService, cfg)
	discHandler := discussionhandler.NewDiscussionHandler(discService)
	expertiseHandler := expertisehandler.NewExpertiseHandler(expertiseService)
	connHandler := connectionhandler.NewConnectionHandler(connService)
	projHandler := projecthandler.NewProjectHandler(projService)
	convHandler := conversationhandler.NewConversationHandler(convService)
	notifHandler := notificationhandler.NewNotificationHandler(notifService)
	dashHandler := dashboardhandler.NewDashboardHandler(dashService)
	homeHandler := homepagehandler.NewHomepageHandler(homeService)
	institutionHandler := institutionhandler.NewInstitutionHandler(institutionService)
	organizationHandler := organizationhandler.NewOrganizationHandler(organizationService)
	locationHandler := locationhandler.NewLocationHandler(locationService)
	industryHandler := industryhandler.NewIndustryHandler(industryService)
	knowledgeHandler := knowledgehandler.NewKnowledgeHandler(knowledgeService)
	r := routes.NewRouter(cfg)
	v1 := r.Group("/api/v1")
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))

	authroutes.RegisterRoutes(v1, protected, authHandler)
	userroutes.RegisterRoutes(v1, protected, userHandler, connHandler.GetSuggestions)
	discussionroutes.RegisterRoutes(v1, protected, discHandler)
	expertiseroutes.RegisterRoutes(v1, expertiseHandler)
	connectionroutes.RegisterRoutes(protected, connHandler)
	projectroutes.RegisterRoutes(v1, protected, projHandler)
	institutionroutes.RegisterRoutes(v1, institutionHandler)
	organizationroutes.RegisterRoutes(v1, organizationHandler)
	locationroutes.RegisterRoutes(v1, locationHandler)
	industryroutes.RegisterRoutes(v1, industryHandler)
	knowledgeroutes.RegisterRoutes(v1, knowledgeHandler)
	conversationroutes.RegisterRoutes(protected, convHandler)
	notificationroutes.RegisterRoutes(protected, notifHandler)
	dashboardroutes.RegisterRoutes(protected, dashHandler)
	homeroutes.RegisterRoutes(v1, homeHandler)
	r.GET("/", func(c *gin.Context) {
		utils.JSONResponse(c, http.StatusOK, "Nalakarsa Backend API Server", gin.H{
			"status":  "healthy",
			"service": "Nalakarsa Backend API Server",
			"version": "1.0.0",
		}, nil)
	})
	log.Printf("Server successfully started on port %s in %s mode\n", cfg.Port, cfg.Env)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Critical Error: Server failed to start: %v", err)
	}
}

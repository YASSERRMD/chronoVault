package main

import (
	"log"
	"os"

	"chronovault/internal/config"
	"chronovault/internal/database"
	"chronovault/internal/handlers"
	"chronovault/internal/middleware"
	"chronovault/internal/repository"
	"chronovault/internal/services"
	"chronovault/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.New(db)
	wsHub := websocket.NewHub()
	go wsHub.Run()

	authService := services.NewAuthService(repo, cfg.JWTSecret)
	contractService := services.NewContractService(repo)
	obligationService := services.NewObligationService(repo, wsHub)
	reportService := services.NewReportService(repo)
	auditService := services.NewAuditService(repo)

	h := handlers.New(authService, contractService, obligationService, reportService, auditService, wsHub)

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleWebSocket(c, wsHub)
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/register", h.Register)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			organizations := protected.Group("/organizations")
			{
				organizations.GET("", h.ListOrganizations)
				organizations.GET("/:id", h.GetOrganization)
				organizations.POST("", middleware.RequireRole("admin"), h.CreateOrganization)
				organizations.PUT("/:id", middleware.RequireRole("admin"), h.UpdateOrganization)
			}

			contracts := protected.Group("/contracts")
			{
				contracts.GET("", h.ListContracts)
				contracts.GET("/:id", h.GetContract)
				contracts.GET("/:id/versions", h.GetContractVersions)
				contracts.GET("/:id/clauses", h.GetContractClauses)
				contracts.POST("", h.CreateContract)
				contracts.PUT("/:id", h.UpdateContract)
				contracts.DELETE("/:id", h.DeleteContract)
			}

			clauses := protected.Group("/clauses")
			{
				clauses.GET("", h.ListClauses)
				clauses.POST("", h.CreateClause)
				clauses.PUT("/:id", h.UpdateClause)
				clauses.DELETE("/:id", h.DeleteClause)
			}

			obligations := protected.Group("/obligations")
			{
				obligations.GET("", h.ListObligations)
				obligations.GET("/:id", h.GetObligation)
				obligations.GET("/:id/history", h.GetObligationHistory)
				obligations.POST("", h.CreateObligation)
				obligations.PUT("/:id", h.UpdateObligation)
				obligations.DELETE("/:id", h.DeleteObligation)
				obligations.POST("/:id/fulfill", h.FulfillObligation)
			}

			reports := protected.Group("/reports")
			{
				reports.GET("/financial-summary", h.GetFinancialSummary)
				reports.GET("/penalty-tracking", h.GetPenaltyTracking)
				reports.GET("/risk-exposure", h.GetRiskExposure)
				reports.GET("/yearly-impact", h.GetYearlyImpact)
			}

			audit := protected.Group("/audit")
			{
				audit.GET("", h.ListAuditLogs)
				audit.GET("/:entity_type/:entity_id", h.GetEntityAuditLogs)
			}

			protected.GET("/me", h.GetCurrentUser)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

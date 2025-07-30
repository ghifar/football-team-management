package main

import (
	"context"
	"football-team-management/cmd/web/handlers"
	"football-team-management/cmd/web/middleware"
	"football-team-management/internal/usecases"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	getUser := usecases.NewGetUserImpl()

	// Initialize auth service with JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "unset-secret"
	}
	authService := usecases.NewAuthService(jwtSecret, getUser)

	pingHandler := handlers.NewPingHandlerImpl()
	userLoginHandler := handlers.NewLoginHandlerImpl(authService)

	var teamRepo usecases.TeamRepository
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	teamRepo = usecases.NewPostgresTeamRepo(pool)
	teamHandler := handlers.NewTeamHandler(teamRepo)

	playerRepo := usecases.NewPostgresPlayerRepo(pool)
	playerHandler := handlers.NewPlayerHandler(playerRepo)

	matchRepo := usecases.NewPostgresMatchRepo(pool)
	matchHandler := handlers.NewMatchHandler(matchRepo)

	matchResultRepo := usecases.NewPostgresMatchResultRepo(pool)
	matchResultHandler := handlers.NewMatchResultHandler(matchResultRepo)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/ping", pingHandler.Handle)
		v1 := api.Group("/v1")
		{
			v1.POST("/login", userLoginHandler.Handle)

			// Protected routes - require JWT authentication
			protected := v1.Group("/")
			protected.Use(middleware.JWTAuth(authService))
			{
				// Team management endpoints - require admin role
				protected.POST("/teams", middleware.RequireRole("admin"), teamHandler.Register)
				protected.PUT("/teams/:name", middleware.RequireRole("admin"), teamHandler.Update)
				protected.DELETE("/teams/:name", middleware.RequireRole("admin"), teamHandler.Delete)
				protected.GET("/teams", teamHandler.List)
				protected.PATCH("/teams/:name/restore", middleware.RequireRole("admin"), teamHandler.Restore)

				// Player management endpoints - require admin role
				protected.POST("/players", middleware.RequireRole("admin"), playerHandler.Register)
				protected.PUT("/players/:playerName", middleware.RequireRole("admin"), playerHandler.Update)
				protected.DELETE("/players/:playerName", middleware.RequireRole("admin"), playerHandler.Delete)
				protected.GET("/players", playerHandler.List)
				protected.GET("/players/team/:teamName", playerHandler.ListByTeam)
				protected.PATCH("/players/:playerName/restore", middleware.RequireRole("admin"), playerHandler.Restore)
				protected.GET("/player/:playerName", playerHandler.GetByName)

				// Match management endpoints - require admin role
				protected.POST("/matches", middleware.RequireRole("admin"), matchHandler.Register)
				protected.PUT("/matches/:id", middleware.RequireRole("admin"), matchHandler.Update)
				protected.DELETE("/matches/:id", middleware.RequireRole("admin"), matchHandler.Delete)
				protected.GET("/matches", matchHandler.List)
				protected.GET("/matches/team/:teamName", matchHandler.ListByTeam)
				protected.PATCH("/matches/:id/restore", middleware.RequireRole("admin"), matchHandler.Restore)
				protected.GET("/match/:id", matchHandler.GetByID)

				// Match result management endpoints - require admin role
				protected.POST("/match-results", middleware.RequireRole("admin"), matchResultHandler.Register)
				protected.PUT("/match-results/:id", middleware.RequireRole("admin"), matchResultHandler.Update)
				protected.DELETE("/match-results/:id", middleware.RequireRole("admin"), matchResultHandler.Delete)
				protected.GET("/match-results", matchResultHandler.List)
				protected.GET("/match-results/match/:matchID", matchResultHandler.GetByMatchID)
				protected.PATCH("/match-results/:id/restore", middleware.RequireRole("admin"), matchResultHandler.Restore)
				protected.GET("/match-result/:id", matchResultHandler.GetByID)
			}
		}
	}
	if err := router.Run(":8080"); err != nil {
		log.Fatal("error initializing server")
	}
}

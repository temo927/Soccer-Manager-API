package http

import (
	"soccer-manager-api/internal/app/auth"
	"soccer-manager-api/internal/app/player"
	"soccer-manager-api/internal/app/team"
	"soccer-manager-api/internal/app/transfer"
	"soccer-manager-api/internal/infrastructure/config"
	"soccer-manager-api/internal/infrastructure/transport/http/handlers"
	"soccer-manager-api/internal/infrastructure/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authUseCase *auth.AuthUseCase,
	teamUseCase *team.TeamUseCase,
	playerUseCase *player.PlayerUseCase,
	transferUseCase *transfer.TransferUseCase,
) *gin.Engine {
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler(authUseCase)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			teamHandler := handlers.NewTeamHandler(teamUseCase)
			teams := protected.Group("/teams")
			{
				teams.GET("/me", teamHandler.GetTeam)
				teams.PUT("/me", teamHandler.UpdateTeam)
				teams.GET("/me/players", teamHandler.GetTeamPlayers)
			}

			playerHandler := handlers.NewPlayerHandler(playerUseCase)
			players := protected.Group("/players")
			{
				players.GET("/:id", playerHandler.GetPlayer)
				players.PUT("/:id", playerHandler.UpdatePlayer)
			}

			transferHandler := handlers.NewTransferHandler(transferUseCase)
			{
				protected.POST("/players/:id/transfer-list", transferHandler.ListPlayer)
				protected.DELETE("/players/:id/transfer-list", transferHandler.RemoveFromTransferList)
				protected.GET("/transfer-list", transferHandler.GetTransferList)
				protected.POST("/transfer-list/:listing_id/buy", transferHandler.BuyPlayer)
			}
		}
	}

	return router
}


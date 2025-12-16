package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"soccer-manager-api/internal/app/auth"
	"soccer-manager-api/internal/app/player"
	"soccer-manager-api/internal/app/team"
	"soccer-manager-api/internal/app/transfer"
	redisCache "soccer-manager-api/internal/infrastructure/cache/redis"
	"soccer-manager-api/internal/infrastructure/config"
	"soccer-manager-api/internal/infrastructure/persistence/postgres"
	httpTransport "soccer-manager-api/internal/infrastructure/transport/http"
	"soccer-manager-api/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Init(cfg.App.Environment); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Logger.Info("Connected to PostgreSQL")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Logger.Info("Connected to Redis")

	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	playerRepo := postgres.NewPlayerRepository(db)
	transferRepo := postgres.NewTransferRepository(db)

	cache := redisCache.NewRedisCache(rdb)

	authUseCase := auth.NewAuthUseCase(
		userRepo,
		teamRepo,
		playerRepo,
		cfg.JWT.Secret,
		cfg.JWT.ExpirationHours,
	)

	teamUseCase := team.NewTeamUseCase(teamRepo, playerRepo, cache)
	playerUseCase := player.NewPlayerUseCase(playerRepo, teamRepo, cache)
	transferUseCase := transfer.NewTransferUseCase(transferRepo, teamRepo, playerRepo, cache)

	router := httpTransport.SetupRouter(
		cfg,
		authUseCase,
		teamUseCase,
		playerUseCase,
		transferUseCase,
	)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Logger.Info("Server starting", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Logger.Info("Server exited")
}

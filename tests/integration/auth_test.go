package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"soccer-manager-api/internal/app/auth"
	"soccer-manager-api/internal/app/player"
	"soccer-manager-api/internal/app/team"
	"soccer-manager-api/internal/app/transfer"
	redisCache "soccer-manager-api/internal/infrastructure/cache/redis"
	"soccer-manager-api/internal/infrastructure/config"
	"soccer-manager-api/internal/infrastructure/persistence/postgres"
	httpTransport "soccer-manager-api/internal/infrastructure/transport/http"
	"soccer-manager-api/tests/testutil"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {

	db, err := testutil.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")


	rdb, err := testutil.SetupTestRedis()
	if err != nil {
		t.Fatalf("Failed to setup test Redis: %v", err)
	}


	userRepo := postgres.NewUserRepository(sqlxDB)
	teamRepo := postgres.NewTeamRepository(sqlxDB)
	playerRepo := postgres.NewPlayerRepository(sqlxDB)
	transferRepo := postgres.NewTransferRepository(sqlxDB)


	cache := redisCache.NewRedisCache(rdb)


	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

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


	gin.SetMode(gin.TestMode)
	router := httpTransport.SetupRouter(
		cfg,
		authUseCase,
		teamUseCase,
		playerUseCase,
		transferUseCase,
	)

	server := httptest.NewServer(router)

	cleanup := func() {
		server.Close()
		testutil.CleanupTestDB(db)
		testutil.CleanupTestRedis(rdb)
	}

	return server, cleanup
}

func TestRegister(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", server.URL+"/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	assert.NotNil(t, result["data"])
}

func TestLogin(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()


	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", server.URL+"/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)


	req, _ = http.NewRequest("POST", server.URL+"/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	assert.NotNil(t, result["data"])
}

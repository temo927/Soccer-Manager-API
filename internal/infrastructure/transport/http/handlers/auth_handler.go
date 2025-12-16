package handlers

import (
	"net/http"

	"soccer-manager-api/internal/app/auth"
	"soccer-manager-api/internal/domain"
	"soccer-manager-api/pkg/localization"
	"soccer-manager-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUseCase *auth.AuthUseCase
}

func NewAuthHandler(authUseCase *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))

	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{err.Error()},
		})
		return
	}

	response, err := h.authUseCase.Register(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrUserAlreadyExists {
			statusCode = http.StatusConflict
			message = localization.GetMessage(lang, "user.already_exists")
			logger.Logger.Warn("Registration failed: user already exists", zap.String("email", req.Email))
		} else {
			logger.Logger.Error("Registration failed", zap.String("email", req.Email), zap.Error(err))
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"errors":  []string{err.Error()},
		})
		return
	}

	logger.Logger.Info("User registered successfully", zap.String("user_id", response.User.ID.String()), zap.String("email", req.Email))

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": localization.GetMessage(lang, "user.created"),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))

	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{err.Error()},
		})
		return
	}

	response, err := h.authUseCase.Login(c.Request.Context(), req)
	if err != nil {
		logger.Logger.Warn("Login failed", zap.String("email", req.Email), zap.Error(err))
		statusCode := http.StatusUnauthorized
		message := localization.GetMessage(lang, "user.invalid_credentials")

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"errors":  []string{err.Error()},
		})
		return
	}

	logger.Logger.Info("User logged in successfully", zap.String("user_id", response.User.ID.String()), zap.String("email", req.Email))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": localization.GetMessage(lang, "user.login.success"),
	})
}


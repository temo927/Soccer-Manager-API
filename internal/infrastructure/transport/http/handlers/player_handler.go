package handlers

import (
	"net/http"

	"soccer-manager-api/internal/app/player"
	"soccer-manager-api/internal/domain"
	"soccer-manager-api/pkg/localization"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	playerUseCase *player.PlayerUseCase
}

func NewPlayerHandler(playerUseCase *player.PlayerUseCase) *PlayerHandler {
	return &PlayerHandler{playerUseCase: playerUseCase}
}

func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	playerID := c.Param("id")

	if _, err := uuid.Parse(playerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{"invalid player ID format"},
		})
		return
	}

	player, err := h.playerUseCase.GetPlayer(c.Request.Context(), playerID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrPlayerNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "player.not_found")
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    player,
	})
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	userID := c.GetString("user_id")
	playerID := c.Param("id")

	if _, err := uuid.Parse(playerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{"invalid player ID format"},
		})
		return
	}

	var req player.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{err.Error()},
		})
		return
	}

	updatedPlayer, err := h.playerUseCase.UpdatePlayer(c.Request.Context(), userID, playerID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrPlayerNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "player.not_found")
		} else if err == domain.ErrPlayerNotOwned {
			statusCode = http.StatusForbidden
			message = localization.GetMessage(lang, "player.not_owned")
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedPlayer,
		"message": localization.GetMessage(lang, "player.updated"),
	})
}


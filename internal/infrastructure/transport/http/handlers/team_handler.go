package handlers

import (
	"net/http"

	"soccer-manager-api/internal/app/team"
	"soccer-manager-api/internal/domain"
	"soccer-manager-api/pkg/localization"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamUseCase *team.TeamUseCase
}

func NewTeamHandler(teamUseCase *team.TeamUseCase) *TeamHandler {
	return &TeamHandler{teamUseCase: teamUseCase}
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	userID := c.GetString("user_id")

	team, err := h.teamUseCase.GetTeam(c.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrTeamNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "team.not_found")
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
		"data":    team,
	})
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	userID := c.GetString("user_id")

	var req team.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{err.Error()},
		})
		return
	}

	updatedTeam, err := h.teamUseCase.UpdateTeam(c.Request.Context(), userID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrTeamNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "team.not_found")
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
		"data":    updatedTeam,
		"message": localization.GetMessage(lang, "team.updated"),
	})
}

func (h *TeamHandler) GetTeamPlayers(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	userID := c.GetString("user_id")

	players, err := h.teamUseCase.GetTeamPlayers(c.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrTeamNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "team.not_found")
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
		"data":    players,
	})
}


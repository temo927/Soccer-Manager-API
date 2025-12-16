package handlers

import (
	"net/http"

	"soccer-manager-api/internal/app/transfer"
	"soccer-manager-api/internal/domain"
	"soccer-manager-api/pkg/localization"
	"soccer-manager-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TransferHandler struct {
	transferUseCase *transfer.TransferUseCase
}

func NewTransferHandler(transferUseCase *transfer.TransferUseCase) *TransferHandler {
	return &TransferHandler{transferUseCase: transferUseCase}
}

func (h *TransferHandler) ListPlayer(c *gin.Context) {
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

	var req transfer.ListPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{err.Error()},
		})
		return
	}

	listing, err := h.transferUseCase.ListPlayer(c.Request.Context(), userID, playerID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrPlayerNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "player.not_found")
		} else if err == domain.ErrPlayerNotOwned {
			statusCode = http.StatusForbidden
			message = localization.GetMessage(lang, "player.not_owned")
		} else if err == domain.ErrPlayerAlreadyListed {
			statusCode = http.StatusConflict
			message = localization.GetMessage(lang, "player.already_listed")
		} else if err == domain.ErrInvalidAskingPrice {
			statusCode = http.StatusBadRequest
			message = localization.GetMessage(lang, "error.validation")
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    listing,
		"message": localization.GetMessage(lang, "player.listed"),
	})
}

func (h *TransferHandler) RemoveFromTransferList(c *gin.Context) {
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

	err := h.transferUseCase.RemoveFromTransferList(c.Request.Context(), userID, playerID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrPlayerNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "player.not_found")
		} else if err == domain.ErrPlayerNotOwned {
			statusCode = http.StatusForbidden
			message = localization.GetMessage(lang, "player.not_owned")
		} else if err == domain.ErrPlayerNotOnTransferList {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "player.not_on_list")
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
		"message": localization.GetMessage(lang, "player.removed_from_list"),
	})
}

func (h *TransferHandler) GetTransferList(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	userID := c.GetString("user_id")

	listings, err := h.transferUseCase.GetTransferList(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.internal"),
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    listings,
	})
}

func (h *TransferHandler) BuyPlayer(c *gin.Context) {
	lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))
	userID := c.GetString("user_id")
	listingID := c.Param("listing_id")

	if _, err := uuid.Parse(listingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": localization.GetMessage(lang, "error.validation"),
			"errors":  []string{"invalid listing ID format"},
		})
		return
	}

	transfer, err := h.transferUseCase.BuyPlayer(c.Request.Context(), userID, listingID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := localization.GetMessage(lang, "error.internal")

		if err == domain.ErrTransferListingNotFound {
			statusCode = http.StatusNotFound
			message = localization.GetMessage(lang, "transfer.listing_not_found")
		} else if err == domain.ErrInsufficientBudget {
			statusCode = http.StatusBadRequest
			message = localization.GetMessage(lang, "transfer.insufficient_budget")
			logger.Logger.Warn("Transfer failed: insufficient budget", zap.String("user_id", userID), zap.String("listing_id", listingID))
		} else if err == domain.ErrTeamFull {
			statusCode = http.StatusBadRequest
			message = localization.GetMessage(lang, "transfer.team_full")
		} else if err == domain.ErrCannotBuyOwnPlayer {
			statusCode = http.StatusBadRequest
			message = localization.GetMessage(lang, "transfer.cannot_buy_own")
		} else {
			logger.Logger.Error("Transfer failed", zap.String("user_id", userID), zap.String("listing_id", listingID), zap.Error(err))
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"errors":  []string{err.Error()},
		})
		return
	}

	logger.Logger.Info("Player transfer completed",
		zap.String("player_id", transfer.PlayerID.String()),
		zap.String("buyer_team_id", transfer.BuyerTeamID.String()),
		zap.String("seller_team_id", transfer.SellerTeamID.String()),
		zap.Float64("transfer_price", transfer.TransferPrice),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transfer,
		"message": localization.GetMessage(lang, "transfer.purchased"),
	})
}


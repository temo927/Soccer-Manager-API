package middleware

import (
	"net/http"
	"strings"

	"soccer-manager-api/pkg/jwt"
	"soccer-manager-api/pkg/localization"
	"soccer-manager-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := localization.GetLanguageFromHeader(c.GetHeader("Accept-Language"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": localization.GetMessage(lang, "error.unauthorized"),
				"errors":  []string{"Authorization header is required"},
			})
			c.Abort()
			return
		}


		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": localization.GetMessage(lang, "error.unauthorized"),
				"errors":  []string{"Invalid authorization header format"},
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwt.ValidateToken(token, jwtSecret)
		if err != nil {
			logger.Logger.Warn("Authentication failed", zap.Error(err), zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": localization.GetMessage(lang, "error.unauthorized"),
				"errors":  []string{err.Error()},
			})
			c.Abort()
			return
		}


		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)

		c.Next()
	}
}


package middleware

import (
	"github.com/gin-gonic/gin"
)



func LocalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}


package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func HeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := viper.GetString("AUTHORIZATION_KEY")
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid 'Authorization' header."})
			return
		}

		if auth != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid 'Authorization' header."})
			return
		}

		c.Next()
	}
}

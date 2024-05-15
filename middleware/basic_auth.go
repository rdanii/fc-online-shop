package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	username := viper.GetString("BASIC_AUTH_USERNAME")
	password := viper.GetString("BASIC_AUTH_PASSWORD")

	return func(c *gin.Context) {
		authUsername, authPassword, ok := c.Request.BasicAuth()

		if !ok || authUsername == "" && authPassword == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Basic Authentication Required",
			})
			return
		}

		if authUsername != username {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Invalid username",
			})
			return
		}

		if authPassword != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Invalid password",
			})
			return
		}

		c.Next()
	}
}

package app

import (
	"net/http"
	"online-shop/delivery"
	"online-shop/repository"
	"online-shop/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(postgresConn *gorm.DB) *gin.Engine {
	r := repository.NewRepository(postgresConn)
	u := usecase.NewUsecase(r)
	d := delivery.NewDelivery(u)

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "route not found",
		})
	})

	api := router.Group("/api/v1")
	api.GET("/products", d.GetProducts)

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

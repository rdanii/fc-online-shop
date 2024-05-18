package app

import (
	"net/http"
	"online-shop/delivery"
	"online-shop/middleware"
	"online-shop/repository"
	"online-shop/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(postgresConn *gorm.DB) *gin.Engine {
	orderRepo := repository.NewOrderRepository(postgresConn)

	r := repository.NewRepository(postgresConn)
	u := usecase.NewUsecase(r, orderRepo)
	d := delivery.NewDelivery(u)

	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderDelivery := delivery.NewOrderDelivery(u, orderUsecase)

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "route not found",
		})
	})

	v1 := router.Group("/api/v1")
	admin := router.Group("/admin")
	admin.Use(middleware.HeaderMiddleware())

	// API Products
	v1.GET("/products", d.GetProducts)
	v1.GET("/products/:id", d.GetProductbyID)
	admin.POST("/products", d.CreateProduct)
	admin.PUT("/products/:id", d.UpdateProduct)
	admin.DELETE("/products/:id", d.DeleteProduct)

	// API Orders
	v1.POST("/checkout", d.Checkout)
	v1.POST("/orders/:id/confirm", orderDelivery.ConfirmOrder)
	v1.GET("/orders/:id", orderDelivery.GetDetailOrder)

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

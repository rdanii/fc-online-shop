package delivery

import (
	"net/http"
	"online-shop/usecase"

	"github.com/gin-gonic/gin"
)

type Delivery interface {
	GetProducts(c *gin.Context)
}

type delivery struct {
	usecase usecase.Usecase
}

func NewDelivery(usecase usecase.Usecase) Delivery {
	return &delivery{usecase}
}

func (d *delivery) GetProducts(c *gin.Context) {
	result, err := d.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

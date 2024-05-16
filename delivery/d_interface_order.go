package delivery

import (
	"net/http"
	"online-shop/model/entity"
	"online-shop/usecase"

	"github.com/gin-gonic/gin"
)

type OrderDelivery interface {
	CreateOrder(c *gin.Context)
}

type orderDelivery struct {
	usecase usecase.Usecase
}

func NewOrderDelivery(usecase usecase.Usecase) OrderDelivery {
	return &orderDelivery{usecase}
}

func (d *orderDelivery) CreateOrder(c *gin.Context) {
	var input entity.Checkout

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, errResult := d.usecase.Checkout(input)
	if errResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errResult.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

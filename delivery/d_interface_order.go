package delivery

import (
	"net/http"
	"online-shop/model/entity"
	"online-shop/usecase"

	"github.com/gin-gonic/gin"
)

type OrderDelivery interface {
	CreateOrder(c *gin.Context)
	ConfirmOrder(c *gin.Context)
	GetDetailOrder(c *gin.Context)
}

type orderDelivery struct {
	usecase      usecase.Usecase
	orderUsecase usecase.OrderUsecase
}

func NewOrderDelivery(usecase usecase.Usecase, orderUsecase usecase.OrderUsecase) OrderDelivery {
	return &orderDelivery{usecase, orderUsecase}
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

func (d *orderDelivery) ConfirmOrder(c *gin.Context) {
	id := c.Param("id")
	var order entity.Confirm

	errBind := c.ShouldBindJSON(&order)
	if errBind != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errBind.Error(),
		})
		return
	}

	result, errResult := d.orderUsecase.Confirm(id, order)
	if errResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errResult.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (d *orderDelivery) GetDetailOrder(c *gin.Context) {
	id := c.Param("id")
	passcode := c.Query("passcode")

	result, err := d.orderUsecase.GetDetailOrder(id, passcode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

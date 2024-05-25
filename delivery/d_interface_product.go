package delivery

import (
	"net/http"
	"online-shop/model/dto"
	"online-shop/model/entity"
	"online-shop/usecase"

	"github.com/gin-gonic/gin"
)

type Delivery interface {
	Checkout(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductbyID(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type delivery struct {
	usecase usecase.Usecase
}

func NewDelivery(usecase usecase.Usecase) Delivery {
	return &delivery{usecase}
}

func (d *delivery) GetProducts(c *gin.Context) {
	result, err := d.usecase.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (d *delivery) GetProductbyID(c *gin.Context) {
	id := c.Param("id")

	result, err := d.usecase.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (d *delivery) CreateProduct(c *gin.Context) {
	var input dto.ReqProduct

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, errResult := d.usecase.Create(c, input)
	if errResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errResult.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (d *delivery) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var input dto.ReqProduct

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, errResult := d.usecase.Update(c, id, input)
	if errResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errResult.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (d *delivery) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := d.usecase.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data successfully deleted",
	})
}

func (d *delivery) Checkout(c *gin.Context) {
	var input entity.Checkout

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, errResult := d.usecase.Checkout(c, input)
	if errResult != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errResult.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

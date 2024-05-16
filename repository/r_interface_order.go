package repository

import (
	"online-shop/model/entity"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order entity.Order, details []entity.OrderDetail) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(order entity.Order, details []entity.OrderDetail) error {
	tx := r.db.Begin()

	errOrder := tx.Create(&order).Error
	if errOrder != nil {
		tx.Rollback()
		return errOrder
	}

	errDetails := tx.Create(&details).Error
	if errDetails != nil {
		tx.Rollback()
		return errDetails
	}

	errCommit := tx.Commit().Error
	if errCommit != nil {
		return errCommit
	}

	return nil
}

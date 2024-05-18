package repository

import (
	"online-shop/model/entity"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order entity.Order, details []entity.OrderDetail) error
	GetByID(id string) (entity.Order, error)
	GetDetailOrders(orderID string) ([]entity.OrderDetail, error)
	Update(order entity.Order) (entity.Order, error)
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
		tx.Rollback()
		return errCommit
	}

	return nil
}

func (r *orderRepository) GetByID(id string) (entity.Order, error) {
	var order entity.Order

	err := r.db.Select("id", "email", "address", "passcode", "grand_total", "paid_at", "paid_bank", "paid_account").Find(&order, "id = ?", id).Error
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *orderRepository) GetDetailOrders(orderID string) ([]entity.OrderDetail, error) {
	var orderDetails []entity.OrderDetail

	err := r.db.Select("id", "order_id", "product_id", "quantity", "price", "total").Find(&orderDetails, "order_id", orderID).Error
	if err != nil {
		return orderDetails, err
	}

	return orderDetails, nil
}

func (r *orderRepository) Update(order entity.Order) (entity.Order, error) {
	err := r.db.Updates(&order).Error
	if err != nil {
		return order, err
	}

	return order, nil
}

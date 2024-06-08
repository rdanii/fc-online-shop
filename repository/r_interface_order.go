package repository

import (
	"context"
	"encoding/json"
	"online-shop/model/entity"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order entity.Order, details []entity.OrderDetail) error
	GetByID(c context.Context, id string) (entity.Order, error)
	GetDetailOrders(c context.Context, orderID string) ([]entity.OrderDetail, error)
	Update(c context.Context, order entity.Order) (entity.Order, error)
}

type orderRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client) OrderRepository {
	return &orderRepository{db, redis}
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

func (r *orderRepository) GetByID(c context.Context, id string) (entity.Order, error) {
	var order entity.Order

	// Cek di Redis
	cachedOrder, err := r.redis.Get(c, id).Result()
	if err == nil {
		json.Unmarshal([]byte(cachedOrder), &order)
		return order, nil
	}

	// Jika tidak ada di Redis, ambil dari database
	rows, err := r.db.Model(&entity.Order{}).
		Select("id", "email", "address", "passcode", "grand_total", "paid_at", "paid_bank", "paid_account").
		Where("id = ?", id).
		Rows()
	if err != nil {
		return order, err
	}
	defer rows.Close()

	if rows.Next() {
		err := r.db.ScanRows(rows, &order)
		if err != nil {
			return order, err
		}
	}

	err = rows.Err()
	if err != nil {
		return order, err
	}

	// Serialisasi data order menjadi JSON
	orderJson, err := json.Marshal(order)
	if err != nil {
		return order, err
	}

	// Menyimpan hasil serialisasi ke Redis
	err = r.redis.Set(c, id, orderJson, 0).Err()
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *orderRepository) GetDetailOrders(c context.Context, orderID string) ([]entity.OrderDetail, error) {
	var orderDetails []entity.OrderDetail

	rows, err := r.db.Model(&entity.OrderDetail{}).
		Select("id", "order_id", "product_id", "quantity", "price", "total").
		Where("order_id = ?", orderID).
		Rows()
	if err != nil {
		return orderDetails, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderDetail entity.OrderDetail
		if err := r.db.ScanRows(rows, &orderDetail); err != nil {
			return orderDetails, err
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	err = rows.Err()
	if err != nil {
		return orderDetails, err
	}

	detailsJson, err := json.Marshal(orderDetails)
	if err != nil {
		return nil, err
	}

	err = r.redis.Set(c, orderID, detailsJson, 0).Err()
	if err != nil {
		return nil, err
	}

	return orderDetails, nil
}

func (r *orderRepository) Update(c context.Context, order entity.Order) (entity.Order, error) {
	err := r.db.Updates(&order).Error
	if err != nil {
		return order, err
	}

	// Hapus kunci-kunci tersebut dari Redis
	err = r.redis.Del(c, order.ID).Err()
	if err != nil {
		return order, err
	}

	return order, nil
}

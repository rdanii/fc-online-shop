package repository

import (
	"context"
	"encoding/json"
	"online-shop/model/entity"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(c context.Context) ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
	Create(product entity.Product) (entity.Product, error)
	Update(product entity.Product) (entity.Product, error)
	Delete(product entity.Product) (entity.Product, error)
}

type repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) Repository {
	return &repository{db, redis}
}

func (r *repository) GetAll(c context.Context) ([]entity.Product, error) {

	var products []entity.Product

	// Mencoba untuk mendapatkan data dari cache Redis
	cachedData, err := r.redis.Get(c, "all_products").Result()
	if err == nil {
		// Jika data ditemukan di cache, kita bisa langsung mengembalikannya
		if err := json.Unmarshal([]byte(cachedData), &products); err != nil {
			return nil, err
		}
		return products, nil
	}

	// Jika data tidak ada di cache, kita harus mengambilnya dari database
	if err := r.db.Select("id", "name", "price").Find(&products, "is_deleted = FALSE").Error; err != nil {
		return nil, err
	}

	// Menyimpan hasil query ke dalam cache Redis
	jsonData, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}
	if err := r.redis.Set(c, "all_products", jsonData, 24*time.Hour).Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) GetByID(id string) (entity.Product, error) {
	var product entity.Product

	err := r.db.Select("id", "name", "price").Find(&product, "is_deleted = FALSE AND id = ?", id).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Create(product entity.Product) (entity.Product, error) {
	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Update(product entity.Product) (entity.Product, error) {
	err := r.db.Updates(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Delete(product entity.Product) (entity.Product, error) {
	err := r.db.Updates(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

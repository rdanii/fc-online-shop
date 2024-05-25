package repository

import (
	"context"
	"encoding/json"
	"online-shop/model/entity"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(c context.Context) ([]entity.Product, error)
	GetByID(c context.Context, id string) (entity.Product, error)
	Create(c context.Context, product entity.Product) (entity.Product, error)
	Update(c context.Context, product entity.Product) (entity.Product, error)
	Delete(c context.Context, product entity.Product) (entity.Product, error)
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
	productKey := viper.GetString("PRODUCTS_KEY")

	// Mencoba untuk mendapatkan data dari cache Redis
	cachedData, err := r.redis.Get(c, productKey).Result()
	if err == nil {
		// Jika data ditemukan di cache, kita bisa langsung mengembalikannya
		if err := json.Unmarshal([]byte(cachedData), &products); err != nil {
			return nil, err
		}
		return products, nil
	}

	// Jika data tidak ada di cache, kita harus mengambilnya dari database
	rows, err := r.db.Model(&products).Select("id, name, price").Where("is_deleted = ?", false).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Menyimpan hasil query ke dalam cache Redis
	jsonData, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	// Simpan data JSON ke Redis
	err = r.redis.Set(c, productKey, jsonData, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) GetByID(c context.Context, id string) (entity.Product, error) {
	var product entity.Product
	productIdKey := viper.GetString("PRODUCT_ID_KEY") + id

	cachedData, err := r.redis.Get(c, productIdKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &product); err != nil {
			return product, err
		}
		return product, nil
	}

	rows, err := r.db.Model(&product).Select("id, name, price").Where("is_deleted = ? AND id = ?", false, id).Rows()
	if err != nil {
		return product, err
	}
	defer rows.Close()

	for rows.Next() {
		err := r.db.ScanRows(rows, &product)
		if err != nil {
			return product, err
		}
	}

	err = rows.Err()
	if err != nil {
		return product, err
	}

	jsonData, err := json.Marshal(product)
	if err != nil {
		return product, err
	}

	err = r.redis.Set(c, productIdKey, jsonData, 24*time.Hour).Err()
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Create(c context.Context, product entity.Product) (entity.Product, error) {
	productKey := viper.GetString("PRODUCTS_KEY")

	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	err = r.redis.Del(c, productKey).Err()
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Update(c context.Context, product entity.Product) (entity.Product, error) {
	productKey := viper.GetString("PRODUCTS_KEY")
	productIdKey := viper.GetString("PRODUCT_ID_KEY") + product.ID

	err := r.db.Updates(&product).Error
	if err != nil {
		return product, err
	}

	// Invalidate the cache for all products and the specific product
	err = r.redis.Del(c, productKey, productIdKey).Err()
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Delete(c context.Context, product entity.Product) (entity.Product, error) {
	productKey := viper.GetString("PRODUCTS_KEY")
	productIdKey := viper.GetString("PRODUCT_ID_KEY") + product.ID

	err := r.db.Updates(&product).Error
	if err != nil {
		return product, err
	}

	// Invalidate the cache for all products and the specific product
	err = r.redis.Del(c, productKey, productIdKey).Err()
	if err != nil {
		return product, err
	}

	return product, nil
}

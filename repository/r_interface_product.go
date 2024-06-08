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
		err := json.Unmarshal([]byte(cachedData), &products)
		if err != nil {
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
		err := r.db.ScanRows(rows, &product)
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
		err := json.Unmarshal([]byte(cachedData), &product)
		if err != nil {
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
	productIdKey := viper.GetString("PRODUCT_ID_KEY") + product.ID

	// Mendapatkan data cache saat ini
	cachedData, err := r.redis.Get(c, productKey).Result()
	if err != nil && err != redis.Nil {
		return product, err
	}

	// Inisialisasi objek untuk menyimpan produk
	var products []entity.Product

	// Jika data ada di cache, kita mengurai JSON
	if cachedData != "" {
		err := json.Unmarshal([]byte(cachedData), &products)
		if err != nil {
			return product, err
		}
	}

	// Menambahkan produk baru ke data cache
	products = append(products, product)

	// Mengupdate cache dengan daftar produk yang baru
	jsonData, err := json.Marshal(products)
	if err != nil {
		return product, err
	}

	// Mengupdate cache dengan daftar produk yang baru
	err = r.redis.Set(c, productKey, jsonData, 24*time.Hour).Err()
	if err != nil {
		return product, err
	}

	// Menyimpan produk sebagai objek tunggal di Redis untuk akses cepat berdasarkan id
	jsonProduct, err := json.Marshal(product)
	if err != nil {
		return product, err
	}

	// Menyimpan produk sebagai objek tunggal di Redis untuk akses cepat berdasarkan id
	err = r.redis.Set(c, productIdKey, jsonProduct, 24*time.Hour).Err()
	if err != nil {
		return product, err
	}

	// Membuat produk di database
	err = r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Update(c context.Context, product entity.Product) (entity.Product, error) {
	productKey := viper.GetString("PRODUCTS_KEY")
	productIdKey := viper.GetString("PRODUCT_ID_KEY") + product.ID

	// Mengupdate cache produk
	jsonData, err := json.Marshal(product)
	if err != nil {
		return product, err
	}

	err = r.redis.Set(c, productIdKey, jsonData, 24*time.Hour).Err()
	if err != nil {
		return product, err
	}

	// Mendapatkan data cache saat ini
	var products []entity.Product
	cachedData, err := r.redis.Get(c, productKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedData), &products)
		if err != nil {
			return product, err
		}

		// Mengupdate produk tertentu dalam daftar cache
		for i, p := range products {
			if p.ID == product.ID {
				products[i] = product
				break
			}
		}

		// Mengupdate cache dengan daftar produk yang telah dimodifikasi
		jsonData, err := json.Marshal(products)
		if err != nil {
			return product, err
		}

		err = r.redis.Set(c, productKey, jsonData, 24*time.Hour).Err()
		if err != nil {
			return product, err
		}
	}

	// Mengupdate produk di database
	err = r.db.Updates(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Delete(c context.Context, product entity.Product) (entity.Product, error) {
	productKey := viper.GetString("PRODUCTS_KEY")
	productIdKey := viper.GetString("PRODUCT_ID_KEY") + product.ID

	// Menghapus produk dari cache produk
	err := r.redis.Del(c, productIdKey).Err()
	if err != nil {
		return product, err
	}

	// Mendapatkan data cache saat ini
	var products []entity.Product
	cachedData, err := r.redis.Get(c, productKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedData), &products)
		if err != nil {
			return product, err
		}

		// Menghapus produk tertentu dari daftar cache
		for i, p := range products {
			if p.ID == product.ID {
				products = append(products[:i], products[i+1:]...)
				break
			}
		}

		// Mengupdate cache dengan daftar produk yang telah dimodifikasi
		jsonData, err := json.Marshal(products)
		if err != nil {
			return product, err
		}

		err = r.redis.Set(c, productKey, jsonData, 24*time.Hour).Err()
		if err != nil {
			return product, err
		}
	}

	// Menghapus produk di database
	err = r.db.Updates(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

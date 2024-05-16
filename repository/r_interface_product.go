package repository

import (
	"online-shop/model/entity"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
	Create(product entity.Product) (entity.Product, error)
	Update(product entity.Product) (entity.Product, error)
	Delete(product entity.Product) (entity.Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll() ([]entity.Product, error) {
	product := []entity.Product{}

	err := r.db.Debug().Select("id", "name", "price").Find(&product, "is_deleted = FALSE").Error
	if err != nil {
		return product, err
	}

	return product, nil
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
	err := r.db.Debug().Updates(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

package repository

import (
	"online-shop/model/entity"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll() ([]entity.Product, error) {
	product := []entity.Product{}

	err := r.db.Select("id", "name", "price").Find(&product, "is_delete = FALSE").Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) GetByID(id string) (entity.Product, error) {
	var product entity.Product

	err := r.db.Select("id", "name", "price").Where("is_deleted = FALSE AND id = ?", id).Find(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

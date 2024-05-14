package repository

import (
	"online-shop/model/entity"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]entity.Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAll() ([]entity.Product, error) {
	data := []entity.Product{}

	err := r.db.Select("id", "name", "price").Find(&data, "is_delete = FALSE").Error
	if err != nil {
		return data, err
	}

	return data, nil
}

package usecase

import (
	"errors"
	"online-shop/model/entity"
	"online-shop/repository"
)

type Usecase interface {
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{repo}
}

func (u *usecase) GetAll() ([]entity.Product, error) {
	result, err := u.repo.GetAll()
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) GetByID(id string) (entity.Product, error) {
	result, err := u.repo.GetByID(id)
	if err != nil {
		return result, err
	}

	if result.ID == "" {
		return result, errors.New("id not found")
	}

	return result, nil
}

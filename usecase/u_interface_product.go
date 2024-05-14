package usecase

import (
	"online-shop/model/entity"
	"online-shop/repository"
)

type Usecase interface {
	GetAll() ([]entity.Product, error)
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

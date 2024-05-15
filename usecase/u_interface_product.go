package usecase

import (
	"errors"
	"online-shop/model/dto"
	"online-shop/model/entity"
	"online-shop/repository"

	"github.com/google/uuid"
)

type Usecase interface {
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
	Create(input dto.ReqProduct) (entity.Product, error)
	Update(inputID dto.RequestID, input dto.ReqProduct) (entity.Product, error)
	Delete(inputID dto.RequestID) error
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
		return result, errors.New("product not found")
	}

	return result, nil
}

func (u *usecase) Create(input dto.ReqProduct) (entity.Product, error) {
	product := entity.Product{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Price:     input.Price,
		IsDeleted: &[]bool{false}[0],
	}

	result, err := u.repo.Create(product)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) Update(productID dto.RequestID, input dto.ReqProduct) (entity.Product, error) {
	product, errProduct := u.repo.GetByID(productID.ID)
	if errProduct != nil {
		return product, errProduct
	}

	if product.ID != productID.ID {
		return product, errors.New("product not found")
	}

	if product.Name == input.Name && product.Price == input.Price {
		return product, errors.New("no changes detected")
	}

	product.Name = input.Name
	product.Price = input.Price

	result, err := u.repo.Update(product)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) Delete(inputID dto.RequestID) error {
	product, err := u.repo.GetByID(inputID.ID)
	if err != nil {
		return err
	}

	if product.ID != inputID.ID {
		return errors.New("product not found")
	}

	product.IsDeleted = &[]bool{true}[0]

	_, errResult := u.repo.Delete(product)
	if errResult != nil {
		return errResult
	}

	return nil
}

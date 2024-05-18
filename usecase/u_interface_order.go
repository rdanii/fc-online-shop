package usecase

import (
	"errors"
	"online-shop/model/entity"
	"online-shop/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type OrderUsecase interface {
	Confirm(id string, input entity.Confirm) (entity.Order, error)
	GetDetailOrder(id string, passcode string) (entity.OrderWithDetail, error)
}

type orderUsecase struct {
	repo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{repo}
}

func (u *orderUsecase) Confirm(id string, input entity.Confirm) (entity.Order, error) {
	order, err := u.repo.GetByID(id)
	if err != nil {
		return order, err
	}

	if order.ID != id {
		return order, errors.New("order not found")
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(input.Passcode))
	if errPass != nil {
		return order, errors.New("invalid passcode")
	}

	if order.GrandTotal != input.Amount {
		return order, errors.New("total amount mismatch: access to orders is not allowed")
	}

	currentTime := time.Now()

	order.Passcode = nil
	order.PaidAt = &currentTime
	order.PaidAccount = &input.AccountNumber
	order.PaidBank = &input.Bank
	order.GrandTotal = input.Amount

	update, errUpdate := u.repo.Update(order)
	if errUpdate != nil {
		return update, errUpdate
	}

	return update, nil
}

func (u *orderUsecase) GetDetailOrder(id string, passcode string) (entity.OrderWithDetail, error) {
	order, err := u.repo.GetByID(id)
	if err != nil {
		return entity.OrderWithDetail{}, err
	}

	if order.ID != id {
		return entity.OrderWithDetail{}, errors.New("order not found")
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode))
	if errPass != nil {
		return entity.OrderWithDetail{}, errors.New("invalid passcode")
	}

	order.Passcode = nil

	orderDetails, errDetails := u.repo.GetDetailOrders(id)
	if errDetails != nil {
		return entity.OrderWithDetail{}, errDetails
	}

	orderWithDetail := entity.OrderWithDetail{
		Order:   order,
		Details: orderDetails,
	}

	return orderWithDetail, nil
}

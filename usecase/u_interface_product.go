package usecase

import (
	"errors"
	"math/rand"
	"online-shop/model/dto"
	"online-shop/model/entity"
	"online-shop/repository"
	"time"

	"github.com/google/uuid"
)

type Usecase interface {
	Checkout(input entity.Checkout) (entity.OrderWithDetail, error)
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
	Create(input dto.ReqProduct) (entity.Product, error)
	Update(inputID dto.RequestID, input dto.ReqProduct) (entity.Product, error)
	Delete(inputID dto.RequestID) error
}

type usecase struct {
	repo      repository.Repository
	orderRepo repository.OrderRepository
}

func NewUsecase(repo repository.Repository, orderRepo repository.OrderRepository) Usecase {
	return &usecase{repo, orderRepo}
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

func (u *usecase) Checkout(input entity.Checkout) (entity.OrderWithDetail, error) {
	// 1. Ambil Produk
	products, err := u.repo.GetAll()
	if err != nil {
		return entity.OrderWithDetail{}, err
	}

	// Membuat pemetaan produk berdasarkan ID untuk akses yang mudah
	productMap := make(map[string]entity.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// 2. Hitung Total Keseluruhan
	var grandTotal int64
	for _, productQty := range input.Products {
		product := productMap[productQty.ID]
		grandTotal += product.Price * int64(productQty.Quantity)
	}

	// 4. Generate Kode Akses
	passcode := generatePasscode(5) // Mengasumsikan Anda memiliki fungsi untuk menghasilkan kode akses

	// 5. Buat Pesanan
	order := entity.Order{
		ID:         uuid.NewString(),
		Email:      input.Email,
		Address:    input.Address,
		GrandTotal: grandTotal,
		Passcode:   &passcode,
	}

	// Mendeklarasikan orderDetails sebagai slice dari entity.OrderDetail
	var orderDetails []entity.OrderDetail

	// Mengisi orderDetails dengan detail pesanan
	for _, productQty := range input.Products {
		product := productMap[productQty.ID]
		orderDetail := entity.OrderDetail{
			ID:        uuid.NewString(),
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  productQty.Quantity,
			Price:     product.Price,
			Total:     product.Price * int64(productQty.Quantity),
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	// 6. Simpan Pesanan dan Detailnya
	err = u.orderRepo.CreateOrder(order, orderDetails)
	if err != nil {
		return entity.OrderWithDetail{}, err
	}

	// Mengembalikan Respon
	orderWithDetail := entity.OrderWithDetail{
		Order:   order,
		Details: orderDetails,
	}

	return orderWithDetail, nil
}

func generatePasscode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passcode := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range passcode {
		passcode[i] = charset[rand.Intn(len(charset))]
	}
	return string(passcode)
}

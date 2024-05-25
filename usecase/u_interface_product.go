package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"online-shop/model/dto"
	"online-shop/model/entity"
	"online-shop/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	Checkout(c context.Context, input entity.Checkout) (entity.OrderWithDetail, error)
	GetAll(c context.Context) ([]entity.Product, error)
	GetByID(c context.Context, id string) (entity.Product, error)
	Create(c context.Context, input dto.ReqProduct) (entity.Product, error)
	Update(c context.Context, id string, input dto.ReqProduct) (entity.Product, error)
	Delete(c context.Context, id string) error
}

type usecase struct {
	repo      repository.Repository
	orderRepo repository.OrderRepository
}

func NewUsecase(repo repository.Repository, orderRepo repository.OrderRepository) Usecase {
	return &usecase{repo, orderRepo}
}

func (u *usecase) GetAll(c context.Context) ([]entity.Product, error) {
	result, err := u.repo.GetAll(c)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) GetByID(c context.Context, id string) (entity.Product, error) {
	result, err := u.repo.GetByID(c, id)
	if err != nil {
		return result, err
	}

	if result.ID == "" {
		return result, errors.New("product not found")
	}

	return result, nil
}

func (u *usecase) Create(c context.Context, input dto.ReqProduct) (entity.Product, error) {
	product := entity.Product{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Price:     input.Price,
		IsDeleted: &[]bool{false}[0],
	}

	result, err := u.repo.Create(c, product)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) Update(c context.Context, id string, input dto.ReqProduct) (entity.Product, error) {
	product, errProduct := u.repo.GetByID(c, id)
	if errProduct != nil {
		return product, errProduct
	}

	if product.ID != id {
		return product, errors.New("product not found")
	}

	if product.Name == input.Name && product.Price == input.Price {
		return product, errors.New("no changes detected")
	}

	product.Name = input.Name
	product.Price = input.Price

	result, err := u.repo.Update(c, product)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *usecase) Delete(c context.Context, id string) error {
	product, err := u.repo.GetByID(c, id)
	if err != nil {
		return err
	}

	if product.ID != id {
		return errors.New("product not found")
	}

	product.IsDeleted = &[]bool{true}[0]

	_, errResult := u.repo.Delete(c, product)
	if errResult != nil {
		return errResult
	}

	return nil
}

func (u *usecase) Checkout(c context.Context, input entity.Checkout) (entity.OrderWithDetail, error) {
	// 1. Ambil Produk dari Repository
	products, err := u.repo.GetAll(c)
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
		product, exists := productMap[productQty.ID]
		if !exists {
			return entity.OrderWithDetail{}, fmt.Errorf("product with ID %s not found", productQty.ID)
		}
		grandTotal += product.Price * int64(productQty.Quantity)
	}

	// 3. Generate Kode Akses
	passcode := generatePasscode(5) // Mengasumsikan Anda memiliki fungsi untuk menghasilkan kode akses

	hashPasscode, errHash := bcrypt.GenerateFromPassword([]byte(passcode), bcrypt.MinCost)
	if errHash != nil {
		return entity.OrderWithDetail{}, errHash
	}

	passHash := string(hashPasscode)

	// 4. Buat Pesanan
	order := entity.Order{
		ID:         uuid.NewString(),
		Email:      input.Email,
		Address:    input.Address,
		GrandTotal: grandTotal,
		Passcode:   &passHash,
	}

	// 5. Buat Detail Pesanan
	var orderDetails []entity.OrderDetail
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

	// 7. Mengembalikan Respon
	orderWithDetail := entity.OrderWithDetail{
		Order:   order,
		Details: orderDetails,
	}

	orderWithDetail.Order.Passcode = &passcode

	return orderWithDetail, nil
}

// Fungsi untuk menghasilkan kode akses
func generatePasscode(length int) string {
	// Charset berisi karakter yang dapat digunakan dalam passcode
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Membuat slice byte dengan panjang sesuai input 'length'
	passcode := make([]byte, length)

	// Seed generator angka acak dengan nilai unik berdasarkan waktu sekarang
	rand.Seed(time.Now().UnixNano())

	// Mengisi setiap indeks slice 'passcode' dengan karakter acak dari 'charset'
	for i := range passcode {
		passcode[i] = charset[rand.Intn(len(charset))]
	}

	// Mengembalikan passcode sebagai string
	return string(passcode)
}

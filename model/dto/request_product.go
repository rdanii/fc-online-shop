package dto

type ReqProduct struct {
	Name  string `json:"name" binding:"required"`
	Price int64  `json:"price" binding:"required"`
}

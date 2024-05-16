package dto

type RequestID struct {
	ID string `uri:"id" binding:"required"`
}
type IDString struct {
	ID string `json:"id" binding:"required"`
}

type ReqProduct struct {
	Name  string `json:"name" binding:"required"`
	Price int64  `json:"price" binding:"required"`
}

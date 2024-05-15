package dto

type RequestID struct {
	ID string `uri:"id" binding:"required"`
}

type ReqProduct struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

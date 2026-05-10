package server

type ProductDTO struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}

func NewProductDTO() ProductDTO {
	return ProductDTO{
		Title:    "",
		Price:    0,
		Quantity: 0,
	}
}

package products

import (
	"errors"

	"github.com/google/uuid"
)

type Product struct {
	UUID     string  `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}

func NewProduct(title string, price float64, quantity uint) (Product, error) {
	if title == "" || price < 0.0 {
		return Product{}, errors.New("Invalid arguments in constructor")
	}
	return Product{
		UUID:     uuid.New().String(),
		Title:    title,
		Price:    price,
		Quantity: quantity,
	}, nil
}

func (p *Product) GetTotalPrice() float64 {
	return p.Price * float64(p.Quantity)
}

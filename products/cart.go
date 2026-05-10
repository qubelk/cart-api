package products

type Cart struct {
	Items      map[string]Product
	TotalPrice float64
}

func NewCart() Cart {
	return Cart{
		Items:      make(map[string]Product),
		TotalPrice: 0.0,
	}
}

func NewCartRef() *Cart {
	return &Cart{
		Items:      make(map[string]Product),
		TotalPrice: 0.0,
	}
}

func (c *Cart) RecalculateTotal() {
	c.TotalPrice = 0.0
	for _, i := range c.Items {
		c.TotalPrice = i.GetTotalPrice()
	}
}

func (c *Cart) AddItem(item Product) {
	if i, exists := c.Items[item.UUID]; exists {
		i.Quantity += item.Quantity
		c.Items[item.UUID] = i
	} else {
		c.Items[item.UUID] = item
	}

	c.RecalculateTotal()
}

func (c *Cart) RemoveItem(id string) error {
	if _, exists := c.Items[id]; !exists {
		return ItemNotExists
	}

	delete(c.Items, id)
	c.RecalculateTotal()
	return nil
}

func (c *Cart) CleanCart() {
	c.Items = make(map[string]Product)
	c.TotalPrice = 0.0
}

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

func (c *Cart) AddItem(item Product) {
	if i, exists := c.Items[item.UUID]; exists {
		i.Quantity += item.Quantity
		c.Items[item.UUID] = i
		c.TotalPrice += i.GetTotalPrice()
		return
	}
	c.Items[item.UUID] = item
}

func (c *Cart) RemoveItem(id string) error {
	i, exists := c.Items[id]
	if !exists {
		return ItemNotExists
	}

	if i.Quantity == 1 {
		c.TotalPrice -= i.GetTotalPrice()
		delete(c.Items, id)
		return nil
	}

	i.Quantity--
	c.Items[id] = i
	c.TotalPrice -= i.Price
	return nil
}

func (c *Cart) CleanCart() {
	for k := range c.Items {
		delete(c.Items, k)
	}
	c.TotalPrice = 0.0
}

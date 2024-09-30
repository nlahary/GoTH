package models

import "time"

type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusAbandoned CartStatus = "abandoned"
	CartStatusCompleted CartStatus = "completed"
)

type CartItem struct {
	ID        int
	ProductID int
	Quantity  int
	Price     float64
	AddedAt   time.Time
}

type Cart struct {
	ID        int
	UserID    int
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    CartStatus
}

func (c *Cart) AddItem(productID, quantity int, price float64) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity += 1
			c.Items[i].Price += price
			return
		}
	}
	c.Items = append(c.Items, CartItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		AddedAt:   time.Now(),
	})
}

func (c *Cart) RemoveItem(productID int) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

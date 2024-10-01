package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

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

type Carts struct {
	db *sql.DB
}

func NewCarts(db *sql.DB) *Carts {
	return &Carts{db: db}
}

func (c *Carts) CreateCart(userID int) error {
	_, err := c.db.Exec("INSERT INTO carts (user_id, status) VALUES (?, ?)", userID, CartStatusActive)
	if err != nil {
		log.Println("Error creating cart:", err)
		return err
	}
	return err
}

func (c *Carts) AddItem(p Product, u Contact, quantity int) error {
	cartId, err := c.GetCartID(u.Id)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("INSERT INTO cartItems (cart_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", cartId, p.Id, quantity, p.Price)
	return err
}

func (c *Carts) GetCartID(userID int) (int, error) {
	var cartID int
	err := c.db.QueryRow("SELECT id FROM carts WHERE user_id = ? AND status = 'active'", userID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("cart not found")
		}
		return 0, err
	}
	return cartID, nil
}

func (c *Carts) RemoveItem(userID, productID int) error {
	cartID, err := c.GetCartID(userID)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("DELETE FROM cartItems WHERE cart_id = ? AND product_id = ?", cartID, productID)
	return err
}

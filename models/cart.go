package models

import (
	"database/sql"
	"log"
	"time"
)

type CartStatus int

const (
	CartStatusActive    CartStatus = 1
	CartStatusAbandoned CartStatus = 2
	CartStatusCompleted CartStatus = 3
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

	_, err := c.db.Exec("INSERT INTO carts (user_id, status_id) VALUES (?, ?)", userID, CartStatusActive)
	if err != nil {
		log.Println("Error creating new cart:", err)
		return err
	}
	return err
}

// Add a product to the cart and return the total number of items in the cart
func (c *Carts) AddItem(p Product, u Contact, quantity int) (int, error) {
	cartId, err := c.GetCartID(u.Id)
	if err != nil {
		return 0, err
	}
	_, err = c.db.Exec("INSERT INTO cartItems (cart_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", cartId, p.Id, quantity, p.Price)
	if err != nil {
		log.Println("Error adding item to cart:", err)
	}
	numItems, err := c.GetNumItems(cartId)
	return numItems, err
}

func (c *Carts) GetCartID(userID int) (int, error) {
	var cartID int
	err := c.db.QueryRow("SELECT id FROM carts WHERE user_id = ? AND status_id = 1", userID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active cart found for user:", userID)
			log.Println("Creating new cart for user:", userID)
			err = c.CreateCart(userID)
			if err != nil {
				return 0, err
			}
			carID, err := c.GetCartID(userID)
			return carID, err
		}
		return 0, err
	}
	return cartID, nil
}

func (c *Carts) GetNumItems(cartId int) (int, error) {

	var numItems int
	err := c.db.QueryRow("SELECT COUNT(*) FROM cartItems WHERE cart_id = ?", cartId).Scan(&numItems)
	if err != nil {
		log.Println("Error getting number of items in cart:", err)
		return 0, err
	}
	return numItems, err
}

func (c *Carts) RemoveItem(userID, productID int) error {
	cartID, err := c.GetCartID(userID)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("DELETE FROM cartItems WHERE cart_id = ? AND product_id = ?", cartID, productID)
	return err
}

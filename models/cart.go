package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
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
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Status    CartStatus `json:"status"`
}

type Carts struct {
	db *sql.DB
}

func NewCarts(db *sql.DB) *Carts {
	return &Carts{db: db}
}

// Sqlite functions

func (c *Carts) CreateCart(userID int) (int, error) {
	result, err := c.db.Exec("INSERT INTO carts (user_id, status_id) VALUES (?, ?)", userID, CartStatusActive)
	if err != nil {
		log.Println("Error creating new cart:", err)
		return 0, err
	}
	cartID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting new cart ID:", err)
		return 0, err
	}
	return int(cartID), nil
}

func (c *Carts) AddItem(p Product, cartId string, quantity int) (int, error) {

	_, err := c.db.Exec("INSERT INTO cartItems (cart_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", cartId, p.Id, quantity, p.Price)
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
			cartId, err := c.CreateCart(userID)
			if err != nil {
				return 0, err
			}
			return cartId, err
		}
		log.Println("Error getting cart ID:", err)
		return 0, err
	}
	return cartID, nil
}

func (c *Carts) GetNumItems(cartId string) (int, error) {

	var numItems int
	err := c.db.QueryRow("SELECT COUNT(*) FROM cartItems WHERE cart_id = ?", cartId).Scan(&numItems)
	if err != nil {
		log.Println("Error getting number of items in cart:", err)
		return 0, err
	}
	return numItems, nil
}

func (c *Carts) RemoveItem(userID, productID int) error {
	cartID, err := c.GetCartID(userID)
	if err != nil {
		return err
	}
	_, err = c.db.Exec("DELETE FROM cartItems WHERE cart_id = ? AND product_id = ?", cartID, productID)
	return err
}

// Redis functions

func GetCartRedis(cartID string, redisClient *redis.Client, ctx context.Context) (map[string]string, error) {
	cartKey := "cart:" + cartID
	cart, err := redisClient.HGetAll(ctx, cartKey).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			log.Println("Cart not found in Redis")
		default:
			log.Println("Error getting cart from Redis:", err)
		}
	}
	return cart, nil
}

func AddToCartRedis(cartID, productID string, quantity int, redisClient *redis.Client, ctx context.Context) error {
	cartKey := "cart:" + cartID
	// Check if item already exists in cart
	item, err := GetItemByIdRedis(cartID, productID, redisClient, ctx)
	if err != nil {
		log.Println("Error getting item by id from Redis:", err)
	}
	newQuantity := item.Quantity + quantity

	err = redisClient.HSet(ctx, cartKey, productID, newQuantity).Err()
	if err != nil {
		log.Println("Error adding item to cart in Redis:", err)
	}
	redisClient.Expire(ctx, cartKey, 24*time.Hour)
	return nil
}

func GetItemsInCartRedis(cartID string, redisClient *redis.Client, ctx context.Context) ([]CartItem, error) {
	cart, err := GetCartRedis(cartID, redisClient, ctx)
	if err != nil {
		return nil, err
	}
	var items []CartItem
	for productID, quantityStr := range cart {
		quantityInt, err := strconv.Atoi(quantityStr)
		if err != nil {
			log.Println("Error converting quantity to int:", err)
		}
		productIDInt, err := strconv.Atoi(productID)
		if err != nil {
			log.Println("Error converting productID to int:", err)
			continue
		}
		items = append(items, CartItem{ProductID: productIDInt, Quantity: quantityInt})
	}
	for _, item := range items {
		println("ItemID: ", item.ProductID, " Quantity: ", item.Quantity)
	}
	return items, nil
}

func GetTotalNbItemsInCartRedis(cartID string, redisClient *redis.Client, ctx context.Context) (int, error) {
	cart, err := GetCartRedis(cartID, redisClient, ctx)
	if err != nil {
		return 0, err
	}
	var totalItems int
	for _, quantityStr := range cart {
		quantityInt, err := strconv.Atoi(quantityStr)
		if err != nil {
			log.Println("Error converting quantity to int:", err)
		}
		totalItems += quantityInt
	}
	log.Printf("Total number of items in cart: %d", totalItems)
	return totalItems, nil
}

func GetItemByIdRedis(cartID, productID string, redisClient *redis.Client, ctx context.Context) (CartItem, error) {
	cart, err := GetCartRedis(cartID, redisClient, ctx)
	if err != nil {
		return CartItem{}, err
	}
	quantityStr, ok := cart[productID]
	if !ok {
		return CartItem{}, nil
	}
	quantityInt, err := strconv.Atoi(quantityStr)
	if err != nil {
		log.Println("Error converting quantity to int:", err)
	}
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		log.Println("Error converting productID to int:", err)
	}
	return CartItem{ProductID: productIDInt, Quantity: quantityInt}, nil
}

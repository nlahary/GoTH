package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	cookies "github.com/nlahary/website/cookies"
	models "github.com/nlahary/website/models"
	templates "github.com/nlahary/website/templates"
)

func HandleCart(
	tmpl *templates.Templates,
	carts *models.Carts,
	products *models.Products,
	users *models.Contacts,
	redisClient *redis.Client,
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:

			cartId := cookies.GetCartCookie(w, r)
			if cartId == "" {
				http.Error(w, "Cart not found", http.StatusNotFound)
				return
			}

			// Case: /cart/{id}
			ProductIdStr := r.URL.Path[len("/cart/"):]
			ProductId, err := strconv.Atoi(ProductIdStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			quantityInt := 1
			models.AddToCartRedis(cartId, ProductIdStr, quantityInt, redisClient, ctx)
			// Get items in cart from Redis
			item, err := models.GetItemByIdRedis(cartId, ProductIdStr, redisClient, ctx)
			if err != nil {
				log.Println("Error getting item by id from Redis:", err)
			}
			log.Printf("Quantity updated for item %d: %d", item.ProductID, item.Quantity)

			_, _ = models.GetItemsInCartRedis(cartId, redisClient, ctx)

			// Get product from SQLite
			product, err := products.GetProductByID(ProductId)
			if err != nil {
				log.Println("Error getting product by id:", err)
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}

			// Add item to cart in SQLite
			_, err = carts.AddItem(*product, cartId, quantityInt)
			if err != nil {
				log.Println("Error adding item to cart:", err)
				http.Error(w, "Error adding item to cart", http.StatusInternalServerError)
				return
			}
			numProducts, _ := models.GetTotalNbItemsInCartRedis(cartId, redisClient, ctx)
			tmpl.ExecuteTemplate(w, "cartCounter", numProducts)
		}

	}
}

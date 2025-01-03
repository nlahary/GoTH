package handlers

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	cookies "github.com/nlahary/website/cookies"
	models "github.com/nlahary/website/models"
	"github.com/nlahary/website/templates"
)

type PageData struct {
	ProductsList []models.Product
	CartCount    int
}

func HandleProductsGet(
	tmpl *templates.Templates,
	products *models.Products,
	redisClient *redis.Client,
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		productsList, err := products.GetAllProducts()
		if err != nil {
			http.Error(w, "Error getting products", http.StatusInternalServerError)
			return
		}
		cartID := cookies.GetCartCookie(w, r)
		nbItems, err := models.GetTotalNbItemsInCartRedis(cartID, redisClient, ctx)
		if err != nil {
			http.Error(w, "Error getting total number of items in cart", http.StatusInternalServerError)
			return
		}
		pageData := PageData{
			ProductsList: productsList,
			CartCount:    nbItems,
		}

		tmpl.ExecuteTemplate(w, "productsPage", pageData)
	}
}

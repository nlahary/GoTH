package handlers

import (
	"context"
	"net/http"

	cookies "github.com/Nathanael-FR/website/cookies"
	models "github.com/Nathanael-FR/website/models"
	"github.com/Nathanael-FR/website/templates"
	"github.com/go-redis/redis/v8"
)

type PageData struct {
	ProductsList []models.Product
	CartCount    int
}

func HandleProducts(
	tmpl *templates.Templates,
	products *models.Products,
	redisClient *redis.Client,
	ctx context.Context,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

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
}

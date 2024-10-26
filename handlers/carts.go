package handlers

import (
	"log"
	"net/http"
	"strconv"

	cookies "github.com/Nathanael-FR/website/cookies"
	models "github.com/Nathanael-FR/website/models"
	templates "github.com/Nathanael-FR/website/templates"
)

func HandleCart(tmpl *templates.Templates, carts *models.Carts, products *models.Products, users *models.Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:

			// Guest user
			cartId := cookies.GetCartCookie(w, r)

			// Case: /cart/{id}
			ProductIdStr := r.URL.Path[len("/cart/"):]
			ProductId, err := strconv.Atoi(ProductIdStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			quantityInt := 1

			product, err := products.GetProductByID(ProductId)
			if err != nil {
				log.Println("Error getting product by id:", err)
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			numProducts, err := carts.AddItem(*product, cartId, quantityInt)
			if err != nil {
				log.Println("Error adding item to cart:", err)
				http.Error(w, "Error adding item to cart", http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "cartCounter", numProducts)
		}

	}
}

package handlers

import (
	"log"
	"net/http"

	models "github.com/Nathanael-FR/website/models"
	"github.com/Nathanael-FR/website/templates"
)

func HandleProducts(tmpl *templates.Templates, products *models.Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			productsList, err := products.GetAllProducts()
			log.Println(productsList)
			if err != nil {
				http.Error(w, "Error getting products", http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "productsPage", productsList)
			tmpl.ExecuteTemplate(w, "cartCounter", 0)
		}
	}
}

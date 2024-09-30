package handlers

import (
	"net/http"

	models "github.com/Nathanael-FR/website/models"
	"github.com/Nathanael-FR/website/templates"
)

func HandleProducts(tmpl *templates.Templates, products *models.Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			products, err := products.GetAllProducts()
			if err != nil {
				http.Error(w, "Error getting products", http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "productsPage", products)
		}
	}
}

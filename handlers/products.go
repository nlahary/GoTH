package handlers

import (
	"net/http"

	models "github.com/Nathanael-FR/website/models"
	"github.com/Nathanael-FR/website/templates"
)

type PageData struct {
	ProductsList []models.Product
	CartCount    int
}

func HandleProducts(tmpl *templates.Templates, products *models.Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			productsList, err := products.GetAllProducts()
			if err != nil {
				http.Error(w, "Error getting products", http.StatusInternalServerError)
				return
			}
			pageData := PageData{
				ProductsList: productsList,
				CartCount:    0,
			}
			tmpl.ExecuteTemplate(w, "productsPage", pageData)
		}
	}
}

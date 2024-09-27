package products

import (
	"net/http"

	"github.com/Nathanael-FR/website/internal/templates"
)

func HandleProducts(tmpl *templates.Templates, products *Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.ExecuteTemplate(w, "displayProducts", products)
		}
	}
}

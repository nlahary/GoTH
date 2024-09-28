package handlers

import (
	"net/http"

	mod "github.com/Nathanael-FR/website/models"
	"github.com/Nathanael-FR/website/templates"
)

func HandleProducts(tmpl *templates.Templates, products *mod.Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.ExecuteTemplate(w, "productsPage", products)
		}
	}
}

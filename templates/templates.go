package templates

import (
	"html/template"
	"log"
	"net/http"

	mod "github.com/nlahary/website/models"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println("Error executing template", name, ":", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func NewTemplates() *Templates {
	// pattern := filepath.Join("templates", "html", "*.html")
	tmpl := template.New("").Funcs(template.FuncMap{
		"batch": mod.BatchProducts,
	})

	// Ensuite, parsez les fichiers
	tmpl = template.Must(tmpl.ParseGlob("templates/html/*.html"))
	return &Templates{
		templates: tmpl,
	}
}

package templates

import (
	"html/template"
	"net/http"

	mod "github.com/Nathanael-FR/website/models"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	t.templates.ExecuteTemplate(w, name, data)
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

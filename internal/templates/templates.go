package templates

import (
	"html/template"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		// template.Must: This is a helper function that ensures
		// the template parsing is successful. If parsing fails, it will panic.
		templates: template.Must(template.ParseGlob("templates/index.html")),
	}
}

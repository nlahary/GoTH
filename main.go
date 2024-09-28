package main

import (
	"log"
	"net/http"

	handlers "github.com/Nathanael-FR/website/handlers"
	md "github.com/Nathanael-FR/website/middlewares"
	mod "github.com/Nathanael-FR/website/models"
	templates "github.com/Nathanael-FR/website/templates"
)

type Count struct {
	Value int
}

func (c *Count) Increment() {
	c.Value++
}

func main() {
	// A ServeMux (short for "Serve Multiplexer") is a type that helps route HTTP requests
	// to the appropriate handler functions based on the URL patterns.
	router := http.NewServeMux()

	//  Templates refer to HTML files that contain placeholders for dynamic content.
	// These placeholders are replaced with actual values when the template is executed.
	tmpl := templates.NewTemplates()

	contacts := &mod.Contacts{}
	count := &Count{Value: 0}
	products := &mod.Products{}

	*contacts = append(*contacts, mod.Contact{ID: 1, Username: "John", Email: "johndoe@gmail.com"})

	// Create 4 products
	*products = append(*products, *mod.CreateProduct("GTX 3090", "fezfezfeze", 399.90))
	*products = append(*products, *mod.CreateProduct("GTX 3080", "fezfezfeze", 299.90))
	*products = append(*products, *mod.CreateProduct("GTX 3070", "fezfezfeze", 199.90))
	*products = append(*products, *mod.CreateProduct("GTX 3060", "fezfezfeze", 99.90))

	// *products = append(*products, mod.Product{Id: 1, Name: "GTX 3090", Desc: "fezfezfeze", Price: 399.90})

	router.Handle("/", md.DetailedLoggingMiddleware(handleIndex(tmpl, count, contacts)))
	router.Handle("/contacts/", md.DetailedLoggingMiddleware(handlers.HandleContacts(tmpl, contacts)))
	router.Handle("/increment", md.DetailedLoggingMiddleware(handleIncrement(tmpl, count)))
	router.Handle("/products/", md.DetailedLoggingMiddleware(handlers.HandleProducts(tmpl, products)))

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	log.Print("Server started on http://localhost:8080")
	log.Fatal(server.ListenAndServe())

}

func handleIndex(tmpl *templates.Templates, count *Count, contacts *mod.Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index", count)
		tmpl.ExecuteTemplate(w, "display", contacts)
	}
}

func handleIncrement(tmpl *templates.Templates, count *Count) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count.Increment()
		tmpl.ExecuteTemplate(w, "count", count)
	}
}

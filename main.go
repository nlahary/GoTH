package main

import (
	"log"
	"net/http"

	ct "github.com/Nathanael-FR/website/internal/contacts"
	handlers "github.com/Nathanael-FR/website/internal/handlers"
	md "github.com/Nathanael-FR/website/internal/middlewares"
	templates "github.com/Nathanael-FR/website/internal/templates"
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

	contacts := &ct.Contacts{}
	count := &Count{Value: 0}

	*contacts = append(*contacts, ct.Contact{ID: 1, Username: "John", Email: "johndoe@gmail.com"})

	router.Handle("/", md.DetailedLoggingMiddleware(handleIndex(tmpl, count, contacts)))
	router.Handle("/contacts/", md.DetailedLoggingMiddleware(handlers.HandleContacts(tmpl, contacts)))
	router.Handle("/increment", md.DetailedLoggingMiddleware(handleIncrement(tmpl, count)))
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	log.Print("Server started on http://localhost:8080")
	log.Fatal(server.ListenAndServe())

}

func handleIndex(tmpl *templates.Templates, count *Count, contacts *ct.Contacts) http.HandlerFunc {
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

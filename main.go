package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	handlers "github.com/Nathanael-FR/website/handlers"
	md "github.com/Nathanael-FR/website/middlewares"
	mod "github.com/Nathanael-FR/website/models"
	templates "github.com/Nathanael-FR/website/templates"
	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Create a context for the application
	var ctx = context.Background()

	// // Create a Redis client
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// Open a database connection
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Print("Database connection established")

	// A ServeMux (short for "Serve Multiplexer") is a type that helps route HTTP requests
	// to the appropriate handler functions based on the URL patterns.
	router := http.NewServeMux()

	//  Templates refer to HTML files that contain placeholders for dynamic content.
	// These placeholders are replaced with actual values when the template is executed.
	tmpl := templates.NewTemplates()

	// Create instances of DB models
	// These instances will be used to interact with the database
	contactsDB := mod.NewContacts(db)
	productsDB := mod.NewProducts(db)
	cartsDB := mod.NewCarts(db)

	router.Handle("/", md.DetailedLoggingMiddleware(handleIndex(tmpl, contactsDB, cartsDB)))
	router.Handle("/contacts/", md.DetailedLoggingMiddleware(handlers.HandleContacts(tmpl, contactsDB)))
	router.Handle("/products/", md.DetailedLoggingMiddleware(handlers.HandleProducts(tmpl, productsDB)))
	router.Handle("/cart/", md.DetailedLoggingMiddleware(handlers.HandleCart(tmpl, cartsDB, productsDB, contactsDB, redisClient, ctx)))
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	log.Print("Server started on http://localhost:8080")
	log.Fatal(server.ListenAndServe())

}

func handleIndex(tmpl *templates.Templates, contacts *mod.Contacts, carts *mod.Carts) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// cartId := cookies.GetCartCookie(w, r)
		// numProducts, _ := carts.GetNumItems(cartId)
		AllContacts, err := contacts.GetAllContacts()
		if err != nil {
			log.Println("Error getting contacts:", err)
			http.Error(w, "Error getting contacts", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "index", AllContacts)
		tmpl.ExecuteTemplate(w, "display", AllContacts)

	}
}

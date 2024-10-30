package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	handlers "github.com/Nathanael-FR/website/handlers"
	kafka "github.com/Nathanael-FR/website/kafka"
	md "github.com/Nathanael-FR/website/middlewares"
	mod "github.com/Nathanael-FR/website/models"
	templates "github.com/Nathanael-FR/website/templates"
	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
)

const (
	kafkaTopic      = "logs"
	kafkaBrokerAddr = "localhost:9092"
	sqliteDB        = "./app.db"
	redisAddr       = "localhost:6379"
)

func main() {

	var ctx = context.Background()

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	db, err := sql.Open("sqlite3", sqliteDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Print("Database connection established")

	router := http.NewServeMux()

	tmpl := templates.NewTemplates()

	contactsDB := mod.NewContacts(db)
	productsDB := mod.NewProducts(db)
	cartsDB := mod.NewCarts(db)

	producer, err := kafka.NewProducer([]string{kafkaBrokerAddr}, kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	router.Handle("/", md.DetailedLoggingMiddleware(handleIndex(tmpl, contactsDB, cartsDB), producer))
	router.Handle("/contacts/", md.DetailedLoggingMiddleware(handlers.HandleContacts(tmpl, contactsDB), producer))
	router.Handle("/products/", md.DetailedLoggingMiddleware(handlers.HandleProducts(tmpl, productsDB, redisClient, ctx), producer))
	router.Handle("/cart/", md.DetailedLoggingMiddleware(handlers.HandleCart(tmpl, cartsDB, productsDB, contactsDB, redisClient, ctx), producer))
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

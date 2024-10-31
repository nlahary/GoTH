package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
	handlers "github.com/nlahary/website/handlers"
	kafka "github.com/nlahary/website/kafka"
	middlewares "github.com/nlahary/website/middlewares"
	mod "github.com/nlahary/website/models"
	templates "github.com/nlahary/website/templates"
)

const (
	kafkaTopic      = "logs"
	sqliteDB        = "./app.db"
	kafkaBrokerAddr = "localhost:9092"
	redisAddr       = "localhost:6379"
	serverAddr      = "localhost:42069"
)

func main() {

	producer, err := kafka.NewProducer([]string{kafkaBrokerAddr}, kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Kafka producer connection established")
	defer producer.Close()

	logger := mod.BasicLogger{
		DefaultLogger: mod.NewLogger(producer, mod.BasicLogSchema),
	}
	httplogger := mod.HttpLogger{
		DefaultLogger: mod.NewLogger(producer, mod.HttpLogSchema),
	}
	var ctx = context.Background()

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	db, err := sql.Open("sqlite3", sqliteDB)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Print("Database connection established")

	router := http.NewServeMux()

	tmpl := templates.NewTemplates()

	contactsDB := mod.NewContacts(db)
	productsDB := mod.NewProducts(db)
	cartsDB := mod.NewCarts(db)

	router.Handle("/", handleIndex(tmpl, contactsDB))
	router.Handle("/contacts/", handlers.HandleContacts(tmpl, contactsDB))
	router.Handle("/products/", handlers.HandleProducts(tmpl, productsDB, redisClient, ctx))
	router.Handle("/cart/", handlers.HandleCart(tmpl, cartsDB, productsDB, contactsDB, redisClient, ctx))
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := http.Server{
		Addr:    serverAddr,
		Handler: middlewares.DetailedLoggingMiddleware(router, &httplogger),
	}
	logger.Print("Server started on http://" + serverAddr)
	logger.Fatal(server.ListenAndServe())

}

func handleIndex(tmpl *templates.Templates, contacts *mod.Contacts) http.HandlerFunc {

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

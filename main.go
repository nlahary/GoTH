package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlahary/website/handlers"
	"github.com/nlahary/website/kafka"
	"github.com/nlahary/website/middlewares"
	"github.com/nlahary/website/models"
	"github.com/nlahary/website/templates"
	"github.com/riferrei/srclient"
)

const (
	CodeExecLogsTopic  = "logs"
	HttpLogsTopic      = "httplogs"
	sqliteDB           = "./app.db"
	kafkaBrokerAddr    = "localhost:9092"
	redisAddr          = "localhost:6379"
	elasticAddr        = "localhost:9200"
	serverAddr         = "localhost:42069"
	schemaRegistryAddr = "localhost:8081"
)

func main() {

	// Add topic schemas to schema registry
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaRegistryAddr)
	err := kafka.RegisterSchemaIfNotExists(schemaRegistryClient, CodeExecLogsTopic, models.BasicLogSchema)
	if err != nil {
		log.Fatal(err)
	}
	err = kafka.RegisterSchemaIfNotExists(schemaRegistryClient, HttpLogsTopic, models.BasicLogSchema)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Schemas registered in schema registry")

	CodeExecLoggerProducer, err := kafka.NewProducer([]string{kafkaBrokerAddr}, CodeExecLogsTopic, schemaRegistryClient)
	if err != nil {
		log.Fatal(err)
	}
	defer CodeExecLoggerProducer.Close()

	HttpLoggerProducer, err := kafka.NewProducer([]string{kafkaBrokerAddr}, HttpLogsTopic, schemaRegistryClient)
	if err != nil {
		log.Fatal(err)
	}
	defer HttpLoggerProducer.Close()

	log.Print("Kafka broker connection established")

	logger := models.CodeExecLogger{Producer: *CodeExecLoggerProducer}
	httplogger := models.HttpLogger{Producer: *HttpLoggerProducer}

	var ctx = context.Background()

	var redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Redis connection established")

	db, err := sql.Open("sqlite3", sqliteDB)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Print("Database connection established")

	router := http.NewServeMux()

	tmpl := templates.NewTemplates()

	contactsDB := models.NewContacts(db)
	productsDB := models.NewProducts(db)
	cartsDB := models.NewCarts(db)

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

func handleIndex(tmpl *templates.Templates, contacts *models.Contacts) http.HandlerFunc {

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

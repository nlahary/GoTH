package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Nathanael-FR/website/middlewares"
)

type Templates struct {
	templates *template.Template
}

// La fonction ExecuteTemplate permet d'exécuter un template
// Elle prend en paramètre un objet http.ResponseWriter, le nom du template à exécuter et les données à passer au template
// Elle utilise l'objet template.Template pour exécuter le template
func (t *Templates) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	t.templates.ExecuteTemplate(w, name, data)
}

// La fonction NewTemplates permet de créer un objet Templates
// Cet objet contient un pointeur vers un objet template.Template
// Cet objet template.Template contient les templates de notre application web
func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

type Contact struct {
	ID       int
	Username string
	Email    string
}

func CreateContact(id int, username, email string) *Contact {
	return &Contact{
		ID:       id,
		Username: username,
		Email:    email,
	}
}

func (c *Contact) Exists(contacts Contacts) bool {
	for _, contact := range contacts {
		if contact.Username == c.Username || contact.Email == c.Email {
			return true
		}
	}
	return false
}

func DeleteContact(id int, contacts *Contacts) error {
	for i, contact := range *contacts {
		if contact.ID == id {
			*contacts = append((*contacts)[:i], (*contacts)[i+1:]...)
			return nil
		}
	}
	return errors.New("contact not found")
}

func GenerateID(contacts Contacts) int {
	return len(contacts) + 1
}

type Contacts = []Contact

type Count struct {
	Value int
}

func (c *Count) Increment() {
	c.Value++
}

type Handler struct{}

func main() {
	// Créer un routeur :
	// Un routeur est un objet qui permet de définir des routes pour notre application web
	// Une route est une association entre une URL et une fonction qui sera exécutée lorsque l'URL est visitée
	router := http.NewServeMux()

	// Créer un objet Templates pour gérer les templates de notre application web :
	// C'est quoi un template ? Un template est un fichier HTML avec des variables à l'intérieur
	// Les variables sont remplacées par des valeurs lors de l'exécution du template
	tmpl := NewTemplates()

	// Initialiser les données de l'application web
	contacts := &Contacts{}
	count := &Count{Value: 0}

	*contacts = append(*contacts, Contact{ID: 1, Username: "John", Email: "johndoe@gmail.com"})

	router.Handle("/", middlewares.DetailedLoggingMiddleware(handleIndex(tmpl, count, contacts)))
	router.Handle("/contacts/", middlewares.DetailedLoggingMiddleware(HandleContacts(tmpl, contacts)))
	router.Handle("/increment", middlewares.DetailedLoggingMiddleware(handleIncrement(tmpl, count)))

	// Créer un serveur web :
	// Un serveur web est un programme qui écoute les requêtes HTTP et y répond
	// Le serveur web utilise un routeur pour déterminer quelle fonction exécuter en fonction de l'URL visitée
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	log.Print("Server started on http://localhost:8080")
	log.Fatal(server.ListenAndServe())

}

func handleIndex(tmpl *Templates, count *Count, contacts *Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index", count)
		tmpl.ExecuteTemplate(w, "display", contacts)
	}
}

func handleIncrement(tmpl *Templates, count *Count) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count.Increment()
		tmpl.ExecuteTemplate(w, "count", count)
	}
}

func HandleContacts(tmpl *Templates, contacts *Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)

		if r.Method == http.MethodPost {
			username := r.FormValue("name")
			email := r.FormValue("email")
			id := GenerateID(*contacts)
			newContact := CreateContact(id, username, email)

			if newContact.Exists(*contacts) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				tmpl.ExecuteTemplate(w, "error", "User already exists")
				return
			}
			log.Println("Contact created:", newContact)
			*contacts = append(*contacts, *newContact)

			// Renvoie uniquement le nouvel utilisateur
			tmpl.ExecuteTemplate(w, "user", newContact)
			return
		}

		if r.Method == http.MethodDelete {
			log.Println("Delete request received")
			idStr := r.URL.Path[len("/contacts/"):]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			err = DeleteContact(id, contacts)
			if err != nil {
				log.Println("Error deleting contact:", err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Println("Contact deleted:", id)
			w.WriteHeader(http.StatusOK) // Change to 200 OK
			return
		}

		log.Printf("Unhandled method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

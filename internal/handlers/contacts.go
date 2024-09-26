package handlers

import (
	"log"
	"net/http"
	"strconv"

	ct "github.com/Nathanael-FR/website/internal/contacts"
	templates "github.com/Nathanael-FR/website/internal/templates"
)

func HandleContacts(tmpl *templates.Templates, contacts *ct.Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)

		// Two cases for a GET request : /contacts/{id} and /contacts/{id}/edit
		if r.Method == http.MethodGet {

			// Case: /contacts/{id}/edit
			if r.URL.Path[len(r.URL.Path)-5:] == "/edit" {
				idStr := r.URL.Path[len("/contacts/") : len(r.URL.Path)-5] // enlever "/edit"
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Println("Error converting id to int:", err)
					http.Error(w, "Invalid ID", http.StatusBadRequest)
					return
				}
				for _, contact := range *contacts {
					if contact.ID == id {
						// Case: /contacts/{id}/edit
						log.Println("Editing contact:", contact)
						// w.Header().Set("Content-Type", "text/html")
						tmpl.ExecuteTemplate(w, "editForm", contact)
						return
					}
				}
				http.Error(w, "Contact not found", http.StatusNotFound)
				return

			} else {
				// Case: /contacts/{id}
				idStr := r.URL.Path[len("/contacts/"):]
				id, err := strconv.Atoi(idStr)
				log.Println("Displaying contact:", id)
				if err != nil {
					log.Println("Error converting id to int:", err)
					http.Error(w, "Invalid ID", http.StatusBadRequest)
					return
				}
				for _, contact := range *contacts {
					if contact.ID == id {
						log.Println("Displaying contact:", contact)
						tmpl.ExecuteTemplate(w, "user", contact)
						return
					}

				}
				http.Error(w, "Contact not found", http.StatusNotFound)
				return
			}

		}

		if r.Method == http.MethodPost {
			username := r.FormValue("name")
			email := r.FormValue("email")
			newContact := ct.CreateContact(username, email)

			if newContact.Exists(*contacts) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				tmpl.ExecuteTemplate(w, "error", "User already exists")
				return
			}
			log.Println("Contact created:", newContact)
			*contacts = append(*contacts, *newContact)

			tmpl.ExecuteTemplate(w, "user", newContact)
			return
		}

		if r.Method == http.MethodDelete {
			idStr := r.URL.Path[len("/contacts/"):]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			err = ct.DeleteContact(id, contacts)
			if err != nil {
				log.Println("Error deleting contact:", err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Println("Contact deleted:", id)
			w.WriteHeader(http.StatusOK) // Change to 200 so that the fetch API doesn't throw an error
			return
		}

		if r.Method == http.MethodPut {
			idStr := r.URL.Path[len("/contacts/"):]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			username := r.FormValue("name")
			email := r.FormValue("email")
			if ct.UsernameOrEmailIsTaken(id, username, email, *contacts) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				tmpl.ExecuteTemplate(w, "error", "User already exists")
				return
			}
			for i, contact := range *contacts {
				if contact.ID == id {
					(*contacts)[i].Update(username, email)
					log.Println("Contact modified.")
					tmpl.ExecuteTemplate(w, "user", (*contacts)[i])
					return
				}
			}
			http.Error(w, "Contact not found", http.StatusNotFound)
			return
		}

		log.Printf("Unhandled method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

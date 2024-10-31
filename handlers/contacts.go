package handlers

import (
	"log"
	"net/http"
	"strconv"

	models "github.com/nlahary/website/models"
	templates "github.com/nlahary/website/templates"
)

func HandleContacts(tmpl *templates.Templates, contacts *models.Contacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Case: /contacts/{id}/edit
			if r.URL.Path[len(r.URL.Path)-5:] == "/edit" {
				idStr := r.URL.Path[len("/contacts/") : len(r.URL.Path)-5]
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Println("Error converting id to int:", err)
					http.Error(w, "Invalid ID", http.StatusBadRequest)
					return
				}

				contact, err := contacts.GetContactByID(id)
				if err != nil {
					log.Println("Error getting contact:", contact)
					http.Error(w, "Contact not found", http.StatusNotFound)
					return
				}
				tmpl.ExecuteTemplate(w, "editForm", contact)

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
				contact, err := contacts.GetContactByID(id)
				if err != nil {
					log.Println("Error getting contact:", contact)
					http.Error(w, "Contact not found", http.StatusNotFound)
					return
				}
				log.Println("Displaying contact:", contact)
				tmpl.ExecuteTemplate(w, "user", contact)
				return
			}

		case http.MethodPost:

			username := r.FormValue("name")
			email := r.FormValue("email")
			newContact := &models.Contact{Username: username, Email: email, Status: models.ContactStatusUser, Uuid: ""}
			userId, err := contacts.InsertContact(newContact)
			if err != nil {
				log.Println("Error inserting contact:", err)
				w.WriteHeader(http.StatusUnprocessableEntity)
				tmpl.ExecuteTemplate(w, "error", "User already exists")
				// http.Error(w, "Error inserting contact", http.StatusInternalServerError)
				return
			}
			newContact, err = contacts.GetContactByID(userId)
			if err != nil {
				log.Println("Error getting contact:", err)
				http.Error(w, "Error getting contact", http.StatusInternalServerError)
				return
			}
			log.Println("Contact created:", newContact)
			tmpl.ExecuteTemplate(w, "user", newContact)
			return

		case http.MethodDelete:
			idStr := r.URL.Path[len("/contacts/"):]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			err = contacts.DeleteContact(id)
			if err != nil {
				log.Println("Error deleting contact:", err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Println("Contact deleted:", id)
			w.WriteHeader(http.StatusOK) // Change to 200 so that the fetch API doesn't throw an error
			return

		case http.MethodPut:

			idStr := r.URL.Path[len("/contacts/"):]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Error converting id to int:", err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			username := r.FormValue("name")
			email := r.FormValue("email")
			contact := models.Contact{Id: id, Username: username, Email: email}
			err = contacts.UpdateContact(&contact)
			if err != nil {
				log.Println("Error updating contact:", err)
				w.WriteHeader(http.StatusUnprocessableEntity)
				tmpl.ExecuteTemplate(w, "error", "User already exists")
				return
			}
			log.Println("Contact modified.")
			tmpl.ExecuteTemplate(w, "user", contact)
			return

		default:

			log.Printf("Unhandled method: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

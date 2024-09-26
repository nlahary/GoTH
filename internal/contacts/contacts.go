package contacts

import (
	"errors"
)

type Contact struct {
	ID       int
	Username string
	Email    string
}

type Contacts = []Contact

var IDcounter = 1

func generateID() int {
	IDcounter++ // Incremente l'IDcounter
	return IDcounter
}

func CreateContact(username, email string) *Contact {
	id := generateID()
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

func UsernameOrEmailIsTaken(id int, username, email string, contacts Contacts) bool {
	for _, contact := range contacts {
		if (contact.Email == email || contact.Username == username) && contact.ID != id {
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

func (c *Contact) Update(username, email string) {
	c.Username = username
	c.Email = email
}

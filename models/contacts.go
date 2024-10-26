package models

import (
	"database/sql"
	"errors"
	"log"
)

type Contact struct {
	Id       int
	Uuid     string
	Username string
	Email    string
	Status   ContactStatus
}

type ContactStatus int

const (
	ContactStatusGuest ContactStatus = 0
	ContactStatusUser  ContactStatus = 1
)

// Contacts is a struct that represents a collection of contacts.
// It has a field db of type *sql.DB, which is a pointer to a sql.DB struct.
// This struct serves as a database gateway for the contacts table.
type Contacts struct {
	db *sql.DB
}

// NewContacts is a constructor function that creates a new Contacts struct.
// It takes a pointer to a sql.DB struct as an argument and returns a pointer to a Contacts struct.
// This function is used to create a new instance of the Contacts struct with the given database connection.
func NewContacts(db *sql.DB) *Contacts {
	return &Contacts{db: db}
}

func (c *Contacts) InsertContact(contact *Contact) (int, error) {
	result, err := c.db.Exec(
		"INSERT INTO contacts (uuid, username, email, status) VALUES (?, ?, ?, ?)",
		contact.Uuid, contact.Username, contact.Email, contact.Status)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Println("New contact inserted with ID:", id)
	return int(id), nil
}

func (c *Contacts) UpdateContact(contact *Contact) error {
	_, err := c.db.Exec("UPDATE contacts SET username = ?, email = ? WHERE id = ?", contact.Username, contact.Email, contact.Id)
	return err
}

func (c *Contacts) DeleteContact(id int) error {
	_, err := c.db.Exec("DELETE FROM contacts WHERE id = ?", id)
	return err
}

func (c *Contacts) GetAllContacts() ([]Contact, error) {
	rows, err := c.db.Query("SELECT id, username, email FROM contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.Id, &contact.Username, &contact.Email); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (c *Contacts) GetContactByID(id int) (*Contact, error) {
	var contact Contact
	err := c.db.QueryRow("SELECT id, username, email FROM contacts WHERE id = ?", id).Scan(&contact.Id, &contact.Username, &contact.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	return &contact, nil
}

func (c *Contacts) GetContactByEmail(email string) (*Contact, error) {
	var contact Contact
	err := c.db.QueryRow("SELECT id, username, email FROM contacts WHERE email = ?", email).Scan(&contact.Id, &contact.Username, &contact.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	return &contact, nil
}

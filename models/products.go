package models

import (
	"database/sql"
	"log"
)

type Product struct {
	Id    int
	Name  string
	Desc  string
	Price float32
}

type Products struct {
	db *sql.DB
}

func NewProducts(db *sql.DB) *Products {
	return &Products{db: db}
}

func (p *Products) InsertProduct(product *Product) error {
	_, err := p.db.Exec("INSERT INTO products (name, desc, price) VALUES (?, ?, ?)", product.Name, product.Desc, product.Price)
	return err
}

func (p *Products) GetAllProducts() ([]Product, error) {
	rows, err := p.db.Query("SELECT id, name, desc, price FROM products")
	if err != nil {
		log.Println("Error getting products:", err)
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Desc, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *Products) GetProductByID(id int) (*Product, error) {
	var product Product
	err := p.db.QueryRow("SELECT id, name, desc, price FROM products WHERE id = ?", id).Scan(&product.Id, &product.Name, &product.Desc, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Products) UpdateProduct(product *Product) error {
	_, err := p.db.Exec("UPDATE products SET name = ?, desc = ?, price = ? WHERE id = ?", product.Name, product.Desc, product.Price, product.Id)
	return err
}

func (p *Products) DeleteProduct(id int) error {
	_, err := p.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}

func BatchProducts(items []Product, size int) [][]Product {
	var batches [][]Product
	for size < len(items) {
		items, batches = items[size:], append(batches, items[0:size])
	}
	batches = append(batches, items)
	return batches
}

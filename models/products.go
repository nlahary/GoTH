package models

type Product struct {
	Id    int
	Name  string
	Desc  string
	Price float32
}

type Products = []Product

var ProductIDcount = 1

func genProductID() int {
	ProductIDcount++
	return ProductIDcount
}

func CreateProduct(name, desc string, price float32) *Product {
	id := genProductID()
	return &Product{
		Id:    id,
		Name:  name,
		Desc:  desc,
		Price: price,
	}
}

func BatchProducts(items []Product, size int) [][]Product {
	var batches [][]Product
	for size < len(items) {
		items, batches = items[size:], append(batches, items[0:size])
	}
	batches = append(batches, items)
	return batches
}

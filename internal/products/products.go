package products

type Product struct {
	Id    int
	Name  string
	Desc  string
	Price float32
}

type Products = []Product

var IDcounter = 1

func generateID() int {
	IDcounter++
	return IDcounter
}

func CreateProduct(name, desc string, price float32) *Product {
	id := generateID()
	return &Product{
		Id:    id,
		Name:  name,
		Desc:  desc,
		Price: price,
	}
}

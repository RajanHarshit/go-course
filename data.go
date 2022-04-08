package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

//Product defines structure of API product.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json: "description"`
	Price       float32 `json:"price"`
	SKU         string  `json: "sku"`
	createdOn   string  `json: "-"`
	updatedOn   string  `json: "-"`
	DeletedOn   string  `json: "-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProduct() Products {
	return productList
}

func UpdateProduct(id int, p *Product) error {
	fmt.Println("Inside UpdateProduct")
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	fmt.Println("p:::", p)
	fmt.Println("pos:::", pos)
	p.ID = id
	productList[pos] = p
	fmt.Println("productList :::", productList)
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	fmt.Println("Inside findProduct")
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.43,
		SKU:         "ABCD",
		createdOn:   time.Now().UTC().String(),
		updatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "pqrs",
		createdOn:   time.Now().UTC().String(),
		updatedOn:   time.Now().UTC().String(),
	},
}

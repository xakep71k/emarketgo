package model

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

type DB struct {
	url string
}

func NewDB(url string) *DB {
	return &DB{url: url}
}

func (db *DB) FindAllProducts(filter *Filter) ([]*Product, error) {
	products, err := loadProducts(db.url)

	if err != nil {
		return nil, err
	}

	if filter.enable != nil {
		products = filterEnabled(*filter.enable, products)
	}

	if filter.sort {
		products = sortProducts(products)
	}
	return products, nil
}

func loadProducts(dataPath string) ([]*Product, error) {
	var products []*Product
	data, err := ioutil.ReadFile(dataPath)

	if err != nil {
		return products, err
	}

	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func filterEnabled(enable bool, products []*Product) []*Product {
	var filtered []*Product

	for _, product := range products {
		if product.Enable == enable {
			filtered = append(filtered, product)
		}
	}

	return filtered
}

func sortProducts(products []*Product) []*Product {
	sort.SliceStable(products, func(i, j int) bool {
		return products[i].Title < products[j].Title
	})

	return products
}

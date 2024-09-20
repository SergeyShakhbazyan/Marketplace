package models

import (
	"github.com/gocql/gocql"
	"time"
)

type Product struct {
	ProductID     gocql.UUID `json:"productID"`
	OwnerID       gocql.UUID `json:"ownerID"`
	Title         string     `json:"title"`
	Images        []string   `json:"images"`
	Description   string     `json:"description"`
	Price         int        `json:"price"`
	CategoryID    gocql.UUID `json:"categoryID"`
	SubcategoryID gocql.UUID `json:"subcategoryID"`
	BrandName     string     `json:"brandName"`
	CreatedAt     time.Time  `json:"createdAt"`
	Views         int        `json:"views"`
	Keywords      []string   `json:"keywords,omitempty"`
}

type ProductFilters struct {
	ProductID     gocql.UUID        `json:"productID"`
	CategoryID    gocql.UUID        `json:"categoryID"`
	SubcategoryID gocql.UUID        `json:"subcategoryID"`
	Filters       map[string]string `json:"filters"`
}

type Filter struct {
	ID    gocql.UUID `json:"id,omitempty"`
	Name  string     `json:"name"`
	Value string     `json:"value"`
}

type ProductWrapContent struct {
	ProductID gocql.UUID `json:"productID"`
	Title     string     `json:"productName"`
	Image     string     `json:"productImage"`
	Price     int        `json:"productPrice"`
}

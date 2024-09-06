package models

import "github.com/gocql/gocql"

type Category struct {
	ID            gocql.UUID    `json:"id"`
	Name          string        `json:"name"`
	Image         string        `json:"image"`
	Subcategories []Subcategory `json:"subcategories"`
}

type Subcategory struct {
	ID         gocql.UUID  `json:"id"`
	ParentID   *gocql.UUID `json:"parentID,omitempty"`
	ParentName string      `json:"parentName,omitempty"`
	Name       string      `json:"name"`
	Brands     []Brand     `json:"brands,omitempty"`
	GroupID    gocql.UUID  `json:"groupID,omitempty"`
	GroupName  string      `json:"groupName,omitempty"`
}

type CategoryWrapContent struct {
	ID    gocql.UUID `json:"id"`
	Name  string     `json:"name"`
	Image string     `json:"image"`
}

type Brand struct {
	ID       gocql.UUID `json:"id"`
	Name     string     `json:"name"`
	ParentID gocql.UUID `json:"parentID,omitempty"`
	Models   []Model    `json:"models,omitempty"`
}

type CategoryInfo struct {
	CategoryName    string     `json:"categoryName"`
	SubcategoryName string     `json:"subcategoryName"`
	CategoryID      gocql.UUID `json:"categoryID"`
	SubcategoryID   gocql.UUID `json:"subcategoryID"`
}

type BrandWrap struct {
	Name         string       `json:"name"`
	ID           gocql.UUID   `json:"id"`
	CategoryInfo CategoryInfo `json:"categoryInfo"`
}

type SubcategoryGroups struct {
	GroupID    gocql.UUID   `json:"groupID"`
	GroupName  string       `json:"groupName"`
	GroupList  []gocql.UUID `json:"groupList"`
	CategoryID gocql.UUID   `json:"categoryID"`
}

type Model struct {
	ID         gocql.UUID          `json:"id"`
	Name       string              `json:"name"`
	ParentID   gocql.UUID          `json:"parentID"`
	Parameters map[string][]string `json:"parameters,omitempty"`
}

type SubcategoryParam struct {
	Subcategory Subcategory `json:"subcategory"`
}

type SubcategoryGroup struct {
	ID            string        `json:"id"`
	GroupName     string        `json:"groupName"`
	Subcategories []Subcategory `json:"subcategories"`
}

type CategoryWithGroups struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Image         string             `json:"image"`
	Subcategories []SubcategoryGroup `json:"subcategories"`
}

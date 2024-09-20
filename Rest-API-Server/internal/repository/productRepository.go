package repository

import (
	"context"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"strings"
	"time"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, product *models.Product, filters *[]map[string]string) error
	DeleteProduct(ctx context.Context, id gocql.UUID) error
	UpdateProduct(ctx context.Context, product models.Product) error
	ProductsWrapsByCategory(ctx context.Context, categoryID gocql.UUID, lastProductID gocql.UUID) ([]models.ProductWrapContent, gocql.UUID, error)
	CreateProductFilters(ctx context.Context, categoryID gocql.UUID, subcategory gocql.UUID, filter models.Filter, productID gocql.UUID) error
	ProductWrapByCategory(ctx context.Context, categoryID gocql.UUID) ([]models.ProductWrapContent, error)
	Products(ctx context.Context) ([]models.ProductWrapContent, error)
	SearchByKeywords(ctx context.Context, searchQuery string) ([]models.ProductWrapContent, error)
	GetProductByOwnerID(ctx context.Context, ownerID gocql.UUID) ([]models.ProductWrapContent, error)
	FindProductsByFilters(ctx context.Context, categoryID gocql.UUID, subcategoryID gocql.UUID, filters models.Filter, limit int) ([]gocql.UUID, error)
	FindProductsByID(ctx context.Context, productID gocql.UUID) (*models.ProductWrapContent, error)
	ProductInfoByID(ctx context.Context, productID gocql.UUID) (*models.Product, *[]models.Filter, error)
}

type productRepository struct {
	session *gocql.Session
}

func NewProductRepository(session *gocql.Session) ProductRepository {
	return &productRepository{session: session}
}

func (r *productRepository) AddProduct(ctx context.Context, product *models.Product, filters *[]map[string]string) error {
	var brandName string
	for _, filterMap := range *filters {
		for filterName, filterValue := range filterMap {
			if filterName == "brand" {
				brandName = filterValue
				break
			}
		}
	}

	query := "INSERT INTO marketplace_keyspace.product(product_id, owner_id, category_id, subcategory_id, title, brandname, description, image, price, keywords, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	if err := r.session.Query(query,
		product.ProductID,
		product.OwnerID,
		product.CategoryID,
		product.SubcategoryID,
		product.Title,
		brandName,
		product.Description,
		product.Images,
		product.Price,
		product.Keywords,
		product.CreatedAt,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}

	query = "UPDATE marketplace_keyspace.product_views SET views = views + 0 WHERE product_id = ?"
	if err := r.session.Query(query, product.ProductID).WithContext(ctx).Exec(); err != nil {
		return err
	}

	for _, filterMap := range *filters {
		for filterName, filterValue := range filterMap {
			query := "INSERT INTO marketplace_keyspace.product_filters(category_id, sub_category_id, filter_name, filter_value, product_id) VALUES (?,?,?,?,?)"
			if err := r.session.Query(query,
				product.CategoryID,
				product.SubcategoryID,
				filterName,
				filterValue,
				product.ProductID,
			).WithContext(ctx).Exec(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, id gocql.UUID) error {
	query := "SELECT category_id, subcategory_id, created_at FROM marketplace_keyspace.product_by_id WHERE product_id = ?"
	var categoryID, subCategoryID gocql.UUID
	var createdAt time.Time
	if err := r.session.Query(query, id).WithContext(ctx).Scan(&categoryID, &subCategoryID, &createdAt); err != nil {
		return err
	}

	filtersQuery := "SELECT filter_name, filter_value FROM marketplace_keyspace.product_filters_by_id WHERE category_id = ? AND sub_category_id = ? AND product_id = ?"
	iter := r.session.Query(filtersQuery, categoryID, subCategoryID, id).WithContext(ctx).Iter()
	defer iter.Close()

	var filterName, filterValue string
	for iter.Scan(&filterName, &filterValue) {
		if err := r.session.Query("DELETE FROM marketplace_keyspace.product_filters WHERE category_id = ? AND sub_category_id = ? AND filter_name = ? AND filter_value = ? AND product_id = ?", categoryID, subCategoryID, filterName, filterValue, id).WithContext(ctx).Exec(); err != nil {
			return err
		}
	}

	if err := iter.Close(); err != nil {
		return err
	}

	if err := r.session.Query("DELETE FROM marketplace_keyspace.product WHERE product_id = ? AND created_at = ? AND category_id = ? AND subcategory_id = ?", id, createdAt, categoryID, subCategoryID).WithContext(ctx).Exec(); err != nil {
		return err
	}

	if err := r.session.Query("DELETE FROM marketplace_keyspace.product_views WHERE product_id = ?", id).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *productRepository) ProductsWrapsByCategory(ctx context.Context, categoryID gocql.UUID, lastProductID gocql.UUID) ([]models.ProductWrapContent, gocql.UUID, error) {
	var productID gocql.UUID
	var productWrapList []models.ProductWrapContent
	pageSize := 2

	var query string
	var iter *gocql.Iter

	if lastProductID == (gocql.UUID{}) {
		query = "SELECT product_id FROM marketplace_keyspace.product WHERE category_id = ? LIMIT ?"
		iter = r.session.Query(query, categoryID, pageSize).WithContext(ctx).Iter()
	} else {
		query = "SELECT product_id FROM marketplace_keyspace.product WHERE category_id = ? AND product_id < ? LIMIT ?"
		iter = r.session.Query(query, categoryID, lastProductID, pageSize).WithContext(ctx).Iter()
	}

	defer iter.Close()

	var productIDList []gocql.UUID
	for iter.Scan(&productID) {
		productIDList = append(productIDList, productID)
	}

	if len(productIDList) == 0 {
		return nil, lastProductID, nil
	}

	query = "SELECT product_id, title, image, price FROM marketplace_keyspace.product WHERE product_id = ?"
	for _, id := range productIDList {
		var product models.ProductWrapContent
		var imageList []string
		err := r.session.Query(query, id).WithContext(ctx).Scan(
			&product.ProductID,
			&product.Title,
			&imageList,
			&product.Price,
		)
		if err != nil {
			return nil, lastProductID, err
		}
		if len(imageList) > 0 {
			product.Image = imageList[0] // Only take the first image
		} else {
			product.Image = ""
		}
		productWrapList = append(productWrapList, product)
	}

	newPagingState := gocql.UUID{}
	if len(productIDList) > 0 {
		newPagingState = productIDList[len(productIDList)-1]
	}

	return productWrapList, newPagingState, nil
}

func (r *productRepository) CreateProductFilters(ctx context.Context, categoryID gocql.UUID, subcategory gocql.UUID, filter models.Filter, productID gocql.UUID) error {
	query := "INSERT INTO marketplace_keyspace.product_filters(category_id, sub_category_id, filter_name, filter_value, product_id) VALUES (?, ?, ?, ?, ?)"
	return r.session.Query(query,
		categoryID,
		subcategory,
		filter.Name,
		filter.Value,
		productID,
	).WithContext(ctx).Exec()
}

func (r *productRepository) ProductWrapByCategory(ctx context.Context, categoryID gocql.UUID) ([]models.ProductWrapContent, error) {
	query := "SELECT product_id, title, image, price FROM marketplace_keyspace.product WHERE category_id = ? limit 8"
	var productWrap models.ProductWrapContent
	var productWrapList []models.ProductWrapContent
	var imageList []string
	iter := r.session.Query(query, categoryID).WithContext(ctx).Iter()
	defer iter.Close()
	for iter.Scan(&productWrap.ProductID, &productWrap.Title, &imageList, &productWrap.Price) {
		if len(imageList) > 0 {
			productWrap.Image = imageList[0]
		} else {
			productWrap.Image = ""
		}
		productWrapList = append(productWrapList, productWrap)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return productWrapList, nil
}

func (r *productRepository) FindProductsByID(ctx context.Context, productID gocql.UUID) (*models.ProductWrapContent, error) {
	query := "SELECT product_id, title, image, price FROM marketplace_keyspace.product_by_id WHERE product_id = ?"
	var productWrap models.ProductWrapContent
	var imageList []string
	if err := r.session.Query(query, productID).WithContext(ctx).Scan(
		&productWrap.ProductID,
		&productWrap.Title,
		&imageList,
		&productWrap.Price,
	); err != nil {
		return nil, err
	}

	if len(imageList) > 0 {
		productWrap.Image = imageList[0]
	} else {
		productWrap.Image = ""
	}

	return &productWrap, nil
}

func (r *productRepository) Products(ctx context.Context) ([]models.ProductWrapContent, error) {
	query := "SELECT product_id, title, image, price FROM marketplace_keyspace.product"
	var productWrap models.ProductWrapContent
	var productWrapList []models.ProductWrapContent
	iter := r.session.Query(query).WithContext(ctx).Iter()
	defer iter.Close()
	var imageList []string
	for iter.Scan(&productWrap.ProductID, &productWrap.Title, &imageList, &productWrap.Price) {
		if len(imageList) > 0 {
			productWrap.Image = imageList[0]
		} else {
			productWrap.Image = ""
		}
		productWrapList = append(productWrapList, productWrap)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return productWrapList, nil
}

//func (r *productRepository) SearchByKeyword(session *gocql.Session, keyword string) ([]models.ProductWrapContent, error) {
//	var products []models.ProductWrapContent
//	normalizedKeyword := strings.ToLower(keyword)
//
//	query := `SELECT product_id, title FROM marketplace_keyspace.product WHERE keywords CONTAINS ?`
//	iter := session.Query(query, normalizedKeyword).Iter()
//
//	var product models.ProductWrapContent
//
//	for iter.Scan(&product.ProductID, &product.Title, &product.Image, &product.Price) {
//		products = append(products, product)
//	}
//	if err := iter.Close(); err != nil {
//		return nil, err
//	}
//
//	return products, nil
//}

func (r *productRepository) SearchByKeywords(ctx context.Context, searchQuery string) ([]models.ProductWrapContent, error) {

	keywords := strings.Fields(searchQuery)
	if len(keywords) == 0 {
		return nil, nil
	}

	productKeywordCount := make(map[gocql.UUID]int)
	productDetails := make(map[gocql.UUID]models.ProductWrapContent)

	for _, keyword := range keywords {
		query := `SELECT product_id, title, image, price FROM marketplace_keyspace.product WHERE keywords CONTAINS ?`
		iter := r.session.Query(query, keyword).Iter()

		var product models.ProductWrapContent
		var imageList []string

		for iter.Scan(&product.ProductID, &product.Title, &imageList, &product.Price) {

			if len(imageList) > 0 {
				product.Image = imageList[0]
			} else {
				product.Image = ""
			}
			if _, exists := productKeywordCount[product.ProductID]; exists {
				productKeywordCount[product.ProductID]++
			} else {
				productKeywordCount[product.ProductID] = 1
				productDetails[product.ProductID] = product
			}
		}
		iter.Close()
	}

	var result []models.ProductWrapContent
	requiredKeywordCount := len(keywords)
	for productID, count := range productKeywordCount {
		if count == requiredKeywordCount {
			result = append(result, productDetails[productID])
		}
	}

	return result, nil
}

func (r *productRepository) GetProductByOwnerID(ctx context.Context, ownerID gocql.UUID) ([]models.ProductWrapContent, error) {
	query := "SELECT product_id, title, image, price FROM marketplace_keyspace.product WHERE owner_id = ?"
	var productWrap models.ProductWrapContent
	var productWrapList []models.ProductWrapContent
	iter := r.session.Query(query, ownerID).WithContext(ctx).Iter()
	defer iter.Close()
	var imageList []string
	for iter.Scan(&productWrap.ProductID, &productWrap.Title, &imageList, &productWrap.Price) {
		if len(imageList) > 0 {
			productWrap.Image = imageList[0]
		} else {
			productWrap.Image = ""
		}
		productWrapList = append(productWrapList, productWrap)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return productWrapList, nil
}

func (r *productRepository) FindProductsByFilters(ctx context.Context, categoryID gocql.UUID, subcategoryID gocql.UUID, filters models.Filter, limit int) ([]gocql.UUID, error) {
	query := "SELECT product_id FROM marketplace_keyspace.product_filters WHERE category_id = ? AND  sub_category_id = ? AND filter_name = ? AND filter_value = ?"
	if limit > 0 {
		query += " LIMIT ?"
	}

	var iter *gocql.Iter
	if limit > 0 {
		iter = r.session.Query(query, categoryID, subcategoryID, filters.Name, filters.Value, limit).WithContext(ctx).Iter()
	} else {
		iter = r.session.Query(query, categoryID, subcategoryID, filters.Name, filters.Value).WithContext(ctx).Iter()
	}

	defer iter.Close()
	var ids []gocql.UUID
	var id gocql.UUID
	for iter.Scan(&id) {
		ids = append(ids, id)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *productRepository) ProductInfoByID(ctx context.Context, productID gocql.UUID) (*models.Product, *[]models.Filter, error) {
	var productInfo models.Product
	productQuery := "SELECT product_id, title, image, description, price, owner_id, created_at, category_id, subcategory_id, brandName FROM marketplace_keyspace.product_by_id WHERE product_id = ?"
	if err := r.session.Query(productQuery, productID).WithContext(ctx).Scan(
		&productInfo.ProductID,
		&productInfo.Title,
		&productInfo.Images,
		&productInfo.Description,
		&productInfo.Price,
		&productInfo.OwnerID,
		&productInfo.CreatedAt,
		&productInfo.CategoryID,
		&productInfo.SubcategoryID,
		&productInfo.BrandName,
	); err != nil {
		return nil, nil, err
	}

	var filters []models.Filter

	productQuery = "SELECT filter_name, filter_value FROM marketplace_keyspace.product_filters WHERE product_id = ?"
	iter := r.session.Query(productQuery, productID).WithContext(ctx).Iter()
	defer iter.Close()

	var filterName, filterValue string

	for iter.Scan(&filterName, &filterValue) {
		filter := models.Filter{
			ID:    gocql.TimeUUID(),
			Name:  filterName,
			Value: filterValue,
		}
		filters = append(filters, filter)
	}

	if err := iter.Close(); err != nil {
		return &productInfo, nil, err
	}

	return &productInfo, &filters, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product models.Product) error {
	query := "UPDATE marketplace_keyspace.product SET owner_id = ?, title = ?, image = ?, description = ?, price = ?,brandName = ?, category_id = ?, subcategory_id = ?,created_at = ?, keywords = ? WHERE product_id = ?"
	return r.session.Query(query,
		product.OwnerID, product.Title, product.Images, product.Description,
		product.Price, product.BrandName, product.CategoryID, product.SubcategoryID,
		product.CreatedAt, product.Views, product.Keywords,
		product.ProductID,
	).WithContext(ctx).Exec()
}

func (r *productRepository) IncrementViews(ctx context.Context, productID gocql.UUID) error {
	query := "UPDATE marketplace_keyspace.product_views SET views = views + 1 WHERE product_id = ?"
	return r.session.Query(query, productID).WithContext(ctx).Exec()
}

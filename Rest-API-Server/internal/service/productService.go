package service

import (
	"context"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) AddProduct(product *models.Product, filters *[]map[string]string) error {
	return s.repo.AddProduct(context.Background(), product, filters)
}

func (s *ProductService) DeleteProduct(productID gocql.UUID) error {
	return s.repo.DeleteProduct(context.Background(), productID)
}

func (s *ProductService) ProductsWrapsByCategory(categoryID gocql.UUID, lastProductID gocql.UUID) ([]models.ProductWrapContent, gocql.UUID, error) {
	return s.repo.ProductsWrapsByCategory(context.Background(), categoryID, lastProductID)
}

func (s *ProductService) ProductInfoByID(productID gocql.UUID) (*models.Product, *[]models.Filter, error) {
	return s.repo.ProductInfoByID(context.Background(), productID)
}

func (s *ProductService) GetProductsByOwnerID(userID gocql.UUID) ([]models.ProductWrapContent, error) {
	return s.repo.GetProductByOwnerID(context.Background(), userID)
}

func (s *ProductService) GetProductsByCategory(categoryID gocql.UUID) ([]models.ProductWrapContent, error) {
	return s.repo.ProductWrapByCategory(context.Background(), categoryID)
}

func (s *ProductService) GetProducts() ([]models.ProductWrapContent, error) {
	return s.repo.Products(context.Background())
}

func (s *ProductService) SearchProducts(searchQuery string) ([]models.ProductWrapContent, error) {
	return s.repo.SearchByKeywords(context.Background(), searchQuery)
}

func (s *ProductService) FindProductsByFilters(categoryID gocql.UUID, subcategoryID gocql.UUID, filter models.Filter, limit int) ([]gocql.UUID, error) {
	return s.repo.FindProductsByFilters(context.Background(), categoryID, subcategoryID, filter, limit)
}

func (s *ProductService) FindProductsByIDs(ctx context.Context, productID gocql.UUID) (*models.ProductWrapContent, error) {
	return s.repo.FindProductsByID(ctx, productID)
}

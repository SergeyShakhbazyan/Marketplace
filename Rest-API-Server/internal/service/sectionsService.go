package service

import (
	"context"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/repository"
)

type SectionsService struct {
	userRepo     repository.UserRepository
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewSectionsService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, userRepo repository.UserRepository) *SectionsService {
	return &SectionsService{userRepo: userRepo, productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *SectionsService) MainCategoriesSection(ctx context.Context) (*models.Section, error) {
	categories, err := s.categoryRepo.ListPopularCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoriesInterface := make([]interface{}, len(categories))
	for i, category := range categories {
		categoriesInterface[i] = category
	}

	section := models.Section{
		SectionID:      gocql.TimeUUID(),
		SectionType:    "categories",
		SectionHeading: "Popular Categories",
		Content:        categoriesInterface,
	}
	return &section, err
}

func (s *SectionsService) MainProductsSections(ctx context.Context) ([]models.Section, error) {
	categories, err := s.categoryRepo.ListPopularCategories(ctx)
	if err != nil {
		return nil, err
	}
	var sections []models.Section

	for _, category := range categories {

		products, err := s.productRepo.ProductWrapByCategory(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		productsInterface := make([]interface{}, len(products))
		for i, product := range products {
			productsInterface[i] = product
		}

		section := models.Section{
			SectionID:      gocql.TimeUUID(),
			SectionType:    "products",
			SectionHeading: category.Name,
			Content:        productsInterface,
		}
		sections = append(sections, section)
	}
	return sections, nil
}

func (s *SectionsService) GetUserProducts(ctx context.Context, ownerID gocql.UUID) (*models.Section, error) {
	products, err := s.productRepo.GetProductByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	productsInterface := make([]interface{}, len(products))
	for i, product := range products {
		productsInterface[i] = product
	}

	section := models.Section{
		SectionID:      gocql.TimeUUID(),
		SectionType:    "products",
		SectionHeading: "User Products",
		Content:        productsInterface,
	}
	return &section, nil
}

func (s *SectionsService) GetProfileInfo(ctx context.Context, userID gocql.UUID) (*models.UserWrapContent, []models.ProductWrapContent, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	products, err := s.productRepo.GetProductByOwnerID(ctx, userID)
	if err != nil {
		return user, nil, err
	}
	return user, products, nil
}

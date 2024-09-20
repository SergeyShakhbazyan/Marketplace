package service

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/repository"
	"marketplace_project/internal/utils"
	"sync"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) AddCategory(category *models.Category) error {
	existingCategory, err := s.repo.GetCategoryDataByName(context.Background(), category.Name)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	if existingCategory != nil {
		return utils.ErrCategoryExists
	}

	return s.repo.CreateCategory(context.Background(), category)
}

func (s *CategoryService) AddSubcategory(subcategoryWithParam *models.SubcategoryParam) error {
	existingSubcategory, err := s.repo.GetSubcategoryDataByName(context.Background(), subcategoryWithParam.Subcategory.Name, subcategoryWithParam.Subcategory.ParentID)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	if existingSubcategory != nil {
		return utils.ErrCategoryExists
	}
	return s.repo.CreateSubcategory(context.Background(), subcategoryWithParam)
}

func (s *CategoryService) GetCategoryDataByName(ctx context.Context, name string) (*models.Category, error) {
	return s.repo.GetCategoryDataByName(ctx, name)
}

func (s *CategoryService) GetSubcategoryDataByName(ctx context.Context, name string, parentID *gocql.UUID) (*models.Subcategory, error) {
	return s.repo.GetSubcategoryDataByName(ctx, name, parentID)
}

func (s *CategoryService) ListMainCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListMainCategories(ctx)
}

func (s *CategoryService) ListSubcategories(ctx context.Context, categoryID gocql.UUID) ([]models.Subcategory, error) {
	return s.repo.ListSubcategoriesByCategoryID(ctx, categoryID)
}

func (s *CategoryService) ListAllCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := s.ListMainCategories(ctx)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	errChan := make(chan error, len(categories))
	subcategoriesMap := make(map[int][]models.Subcategory)
	var mu sync.Mutex

	for i, category := range categories {
		wg.Add(1)
		go func(index int, categoryID gocql.UUID) {
			defer wg.Done()
			subcategories, err := s.ListSubcategories(ctx, categoryID)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			subcategoriesMap[index] = subcategories
			mu.Unlock()
		}(i, category.ID)
	}

	wg.Wait()
	close(errChan)

	var firstError error
	for err := range errChan {
		if firstError == nil {
			firstError = err
		}
	}

	if firstError != nil {
		return nil, err
	}

	for i, subcategories := range subcategoriesMap {
		categories[i].Subcategories = subcategories
	}
	return categories, nil
}

func (s *CategoryService) PopularCategories(ctx context.Context) ([]models.CategoryWrapContent, error) {
	return s.repo.ListPopularCategories(ctx)
}

func (s *CategoryService) FieldsBySubcategory(ctx context.Context, subcategoryID gocql.UUID) (map[string]interface{}, error) {
	return s.repo.FieldsBySubcategory(ctx, subcategoryID)
}

func (s *CategoryService) BrandsOfSubcategory(ctx context.Context, subcategoryID gocql.UUID) ([]models.Brand, error) {
	return s.repo.BrandsBySubcategory(ctx, subcategoryID)
}

func (s *CategoryService) ModelsByBrands(ctx context.Context, brandID gocql.UUID) ([]models.Model, error) {
	return s.repo.ModelsByBrands(ctx, brandID)
}

func (s *CategoryService) ParametersOfModels(ctx context.Context, modelID gocql.UUID) (map[string][]string, error) {
	return s.repo.ParametersOfModels(ctx, modelID)
}

func (s *CategoryService) SearchByBrands(ctx context.Context, search string) ([]models.BrandWrap, error) {
	return s.repo.SearchByBrands(ctx, search)
}

func (s *CategoryService) SearchByKeywords(ctx context.Context, search string) ([]string, error) {
	return s.repo.SearchByKeywords(ctx, search)
}

func (s *CategoryService) GetBrandsByProductName(ctx context.Context, search string) ([]models.BrandWrap, error) {
	return s.repo.GetBrandsByProductName(ctx, search)
}

func (s *CategoryService) InsertSubcategoriesToGroup(ctx context.Context, group *models.SubcategoryGroups) error {
	return s.repo.InsertSubcategoriesToGroup(ctx, group)
}

func (s *CategoryService) GetModelsByBrandName(ctx context.Context, brandName string) ([]string, error) {
	return s.repo.GetModelsByBrandName(ctx, brandName)
}

//func (s *CategoryService) ListMainCategories(ctx context.Context) ([]models.Category, error) {
//	return s.repo.ListMainCategories(ctx)
//}
//
//func (s *CategoryService) ListSubcategories(ctx context.Context, parentID gocql.UUID) ([]models.Category, error) {
//	return s.repo.ListSubCategoriesByParentID(ctx, parentID)
//}

//func (s *CategoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
//	categories, err := s.repo.FetchAllCategories(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	return s.buildCategoryHierarchy(categories), nil
//}

//func (s *CategoryService) buildCategoryHierarchy(categories []models.Category) []models.Category {
//	lookup := make(map[gocql.UUID]*models.Category)
//	var rootCategories []models.Category
//
//	for i := range categories {
//		cat := &categories[i]
//		lookup[cat.ID] = cat
//	}
//
//	for i := range categories {
//		cat := &categories[i]
//		if cat.ParentID != nil && *cat.ParentID != (gocql.UUID{}) {
//			parent, exists := lookup[*cat.ParentID]
//			if exists {
//				if parent.Subcategories == nil {
//					parent.Subcategories = []models.Category{}
//				}
//				parent.Subcategories = append(parent.Subcategories, *cat)
//			}
//		} else {
//			rootCategories = append(rootCategories, *cat)
//		}
//	}
//
//	return rootCategories
//}

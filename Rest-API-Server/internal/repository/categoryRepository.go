package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/utils"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	CreateSubcategory(ctx context.Context, subcategoryWithParam *models.SubcategoryParam) error
	GetCategoryDataByName(ctx context.Context, name string) (*models.Category, error)
	GetSubcategoryDataByName(ctx context.Context, name string, parentID *gocql.UUID) (*models.Subcategory, error)
	ListMainCategories(ctx context.Context) ([]models.Category, error)
	ListSubcategoriesByCategoryID(ctx context.Context, categoryID gocql.UUID) ([]models.Subcategory, error)
	FetchAllCategories(ctx context.Context) ([]models.Category, error)
	DeleteCategory(ctx context.Context, id gocql.UUID) error
	ListPopularCategories(ctx context.Context) ([]models.CategoryWrapContent, error)
	FieldsBySubcategory(ctx context.Context, subcategoryID gocql.UUID) (map[string]interface{}, error)
	BrandsBySubcategory(ctx context.Context, subcategoryID gocql.UUID) ([]models.Brand, error)
	ModelsByBrands(ctx context.Context, brandID gocql.UUID) ([]models.Model, error)
	ParametersOfModels(ctx context.Context, modelID gocql.UUID) (map[string][]string, error)
	SearchByBrands(ctx context.Context, search string) ([]models.BrandWrap, error)
	SearchByKeywords(ctx context.Context, search string) ([]string, error)
	GetBrandsByProductName(ctx context.Context, search string) ([]models.BrandWrap, error)
	GetModelsByBrandName(ctx context.Context, brandName string) ([]string, error)
	InsertSubcategoriesToGroup(ctx context.Context, group *models.SubcategoryGroups) error
}

type categoryRepository struct {
	session *gocql.Session
}

func NewCategoryRepository(session *gocql.Session) CategoryRepository {
	return &categoryRepository{session: session}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := "INSERT INTO marketplace_keyspace.category(id, name, image) VALUES (?, ?, ?)"
	return r.session.Query(query,
		category.ID,
		category.Name,
		category.Image,
	).WithContext(ctx).Exec()
}
func (r *categoryRepository) CreateSubcategory(ctx context.Context, subcategoryWithParam *models.SubcategoryParam) error {
	//filtersAndInputsJSON, err := json.Marshal(subcategory.FiltersAndInputs)
	//if err != nil {
	//	return err
	//}

	query := "INSERT INTO marketplace_keyspace.subcategories(subcategory_id, category_id, name) VALUES (?, ?, ?)"
	err := r.session.Query(query,
		subcategoryWithParam.Subcategory.ID,
		subcategoryWithParam.Subcategory.ParentID,
		subcategoryWithParam.Subcategory.Name,
	).WithContext(ctx).Exec()
	if err != nil {
		return err
	}

	for _, brand := range subcategoryWithParam.Subcategory.Brands {
		brand.ID = gocql.TimeUUID()
		brandQuery := "INSERT INTO marketplace_keyspace.brands (subcategory_id, id, name) VALUES (?, ?, ?)"
		if err := r.session.Query(brandQuery, subcategoryWithParam.Subcategory.ID, brand.ID, brand.Name).WithContext(ctx).Exec(); err != nil {
			return err
		}
		for _, model := range brand.Models {
			model.ID = gocql.TimeUUID()
			modelQuery := "INSERT INTO marketplace_keyspace.models (id, brand_id, name) VALUES (?, ?, ?)"
			if err := r.session.Query(modelQuery, model.ID, brand.ID, model.Name).WithContext(ctx).Exec(); err != nil {
				return err
			}

			for paramName, paramValues := range model.Parameters {
				paramQuery := "INSERT INTO marketplace_keyspace.model_parameters (modelID, parameterName, parameterValue) VALUES (?, ?, ?)"
				if err := r.session.Query(paramQuery, model.ID, paramName, paramValues).WithContext(ctx).Exec(); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func (r *categoryRepository) InsertSubcategoriesToGroup(ctx context.Context, group *models.SubcategoryGroups) error {
	if group.GroupID == (gocql.UUID{}) {
		group.GroupID = gocql.TimeUUID()
	}
	query := "INSERT INTO marketplace_keyspace.subcategorygroups(categoryid, groupname, groupid, subcategory_ids) VALUES (?, ?, ?, ?)"
	if err := r.session.Query(query, group.CategoryID, group.GroupName, group.GroupID, group.GroupList).WithContext(ctx).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) GetCategoryDataByName(ctx context.Context, name string) (*models.Category, error) {
	query := "SELECT id, name, image FROM marketplace_keyspace.category WHERE name = ?"
	var category models.Category
	err := r.session.Query(query, name).WithContext(ctx).Scan(&category.ID, &category.Name, &category.Image)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetSubcategoryDataByName(ctx context.Context, name string, parentID *gocql.UUID) (*models.Subcategory, error) {
	var subcategory models.Subcategory
	query := "SELECT subcategory_id,category_id,name FROM marketplace_keyspace.subcategories_by_name_category_id WHERE category_id = ? AND name = ?"
	err := r.session.Query(query, *parentID, name).WithContext(ctx).Scan(&subcategory.ID, &subcategory.ParentID, &subcategory.Name)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &subcategory, nil
}

func (r *categoryRepository) ListMainCategories(ctx context.Context) ([]models.Category, error) {
	var category models.Category
	var categoryList []models.Category

	query := "SELECT id, name, image FROM marketplace_keyspace.category"
	iter := r.session.Query(query).WithContext(ctx).Iter()
	defer iter.Close()
	for iter.Scan(&category.ID, &category.Name, &category.Image) {
		categoryList = append(categoryList, category)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return categoryList, nil
}

func (r *categoryRepository) ListSubcategoriesByCategoryID(ctx context.Context, categoryID gocql.UUID) ([]models.Subcategory, error) {
	var subcategory models.Subcategory
	var subcategoriesList []models.Subcategory

	query := "SELECT subcategory_id, category_id, name, groupName FROM marketplace_keyspace.subcategories WHERE category_id = ?"
	iter := r.session.Query(query, categoryID).WithContext(ctx).Iter()
	defer iter.Close()
	for iter.Scan(&subcategory.ID, &subcategory.ParentID, &subcategory.Name, &subcategory.GroupName) {
		subcategoriesList = append(subcategoriesList, subcategory)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return subcategoriesList, nil
}

func (r *categoryRepository) ListPopularCategories(ctx context.Context) ([]models.CategoryWrapContent, error) {
	var category models.CategoryWrapContent
	var categoryList []models.CategoryWrapContent

	query := "SELECT id, name, image FROM marketplace_keyspace.category LIMIT 10"
	iter := r.session.Query(query).WithContext(ctx).Iter()
	defer iter.Close()
	for iter.Scan(&category.ID, &category.Name, &category.Image) {
		categoryList = append(categoryList, category)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return categoryList, nil
}

//func (r *categoryRepository) ListMainCategories(ctx context.Context) ([]models.Category, error) {
//	query := "SELECT id, parent_id, name, image FROM marketplace_keyspace.category WHERE parent_id = ?"
//	var category models.Category
//	var categoriesList []models.Category
//	iter := r.session.Query(query, gocql.UUID{}).WithContext(ctx).Iter()
//	defer iter.Close()
//	for iter.Scan(&category.ID, &category.Name, &category.Image) {
//		categoriesList = append(categoriesList, category)
//	}
//	if err := iter.Close(); err != nil {
//		return nil, err
//	}
//
//	return categoriesList, nil
//}
//
//func (r *categoryRepository) ListSubCategoriesByParentID(ctx context.Context, parentID gocql.UUID) ([]models.Category, error) {
//	query := "SELECT id, name, image FROM marketplace_keyspace.category WHERE parent_id = ?"
//	var category models.Category
//	var categoriesList []models.Category
//	iter := r.session.Query(query, parentID).WithContext(ctx).Iter()
//	defer iter.Close()
//	for iter.Scan(&category.ID, &category.ParentID, &category.Name, &category.Image) {
//		categoriesList = append(categoriesList, category)
//	}
//	if err := iter.Close(); err != nil {
//		return nil, err
//	}
//
//	return categoriesList, nil
//}

func (r *categoryRepository) FetchAllCategories(ctx context.Context) ([]models.Category, error) {
	query := "SELECT id, name, image FROM marketplace_keyspace.category"
	var category models.Category
	var categoriesList []models.Category
	iter := r.session.Query(query).WithContext(ctx).Iter()
	defer iter.Close()
	for iter.Scan(&category.ID, &category.Name, &category.Image) {
		categoriesList = append(categoriesList, category)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return categoriesList, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id gocql.UUID) error {
	query := "DELETE FROM marketplace_keyspace.category WHERE id = ?"
	return r.session.Query(query, id).WithContext(ctx).Exec()
}

func (r *categoryRepository) BrandsBySubcategory(ctx context.Context, subcategoryID gocql.UUID) ([]models.Brand, error) {
	var brands []models.Brand
	query := "SELECT id, name FROM marketplace_keyspace.brands WHERE subcategory_id = ?"
	iter := r.session.Query(query, subcategoryID).WithContext(ctx).Iter()

	var brandID gocql.UUID
	var brandName string

	for iter.Scan(&brandID, &brandName) {
		brands = append(brands, models.Brand{
			ID:       brandID,
			Name:     brandName,
			ParentID: subcategoryID,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *categoryRepository) ModelsByBrands(ctx context.Context, brandID gocql.UUID) ([]models.Model, error) {
	var modelsList []models.Model

	query := "SELECT id, name FROM marketplace_keyspace.models WHERE brand_id = ?"
	iter := r.session.Query(query, brandID).WithContext(ctx).Iter()

	var modelID gocql.UUID
	var modelName string

	for iter.Scan(&modelID, &modelName) {
		modelsList = append(modelsList, models.Model{
			ID:       modelID,
			Name:     modelName,
			ParentID: brandID,
			//Parameters: parameters,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return modelsList, nil
}

func (r *categoryRepository) ParametersOfModels(ctx context.Context, modelID gocql.UUID) (map[string][]string, error) {
	parameters := make(map[string][]string) // Initialize the map
	var paramName string
	var paramValue []string

	query := "SELECT parameterName, parameterValue FROM marketplace_keyspace.model_parameters WHERE modelID = ?"
	iter := r.session.Query(query, modelID).WithContext(ctx).Iter()

	for iter.Scan(&paramName, &paramValue) {
		parameters[paramName] = paramValue
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return parameters, nil
}

type Category struct {
	name string
	id   gocql.UUID
}

func (r *categoryRepository) SearchByBrands(ctx context.Context, search string) ([]models.BrandWrap, error) {
	var name string
	var id gocql.UUID
	var subcategoryID gocql.UUID
	var brandList []models.BrandWrap

	query := "SELECT name, id, subcategory_id FROM marketplace_keyspace.brands_by_name WHERE name LIKE ? ALLOW FILTERING "
	iter := r.session.Query(query, search+"%").WithContext(ctx).Iter()

	for iter.Scan(&name, &id, &subcategoryID) {
		brand := models.BrandWrap{
			Name: name,
			ID:   id,
		}
		var categoryInfo models.CategoryInfo
		subcategoryQuery := "SELECT name, category_id FROM marketplace_keyspace.subcategories WHERE subcategory_id = ?"
		subcategoryIter := r.session.Query(subcategoryQuery, subcategoryID).WithContext(ctx).Iter()

		for subcategoryIter.Scan(&categoryInfo.SubcategoryName, &categoryInfo.CategoryID) {
			categoryQuery := "SELECT name FROM marketplace_keyspace.category WHERE id = ?"
			if err := r.session.Query(categoryQuery, categoryInfo.CategoryID).WithContext(ctx).Scan(&categoryInfo.CategoryName); err != nil {
				return nil, err
			}
		}
		if err := subcategoryIter.Close(); err != nil {
			return nil, err
		}
		categoryInfo.SubcategoryID = subcategoryID
		brand.CategoryInfo = categoryInfo
		brandList = append(brandList, brand)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return brandList, nil
}

func (r *categoryRepository) SearchByKeywords(ctx context.Context, search string) ([]string, error) {
	var name string
	var listHelpingNames []string
	query := "SELECT name FROM marketplace_keyspace.models WHERE name LIKE ? limit 5 ALLOW FILTERING"
	iter := r.session.Query(query, search+"%").WithContext(ctx).Iter()
	for iter.Scan(&name) {
		listHelpingNames = append(listHelpingNames, name)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return listHelpingNames, nil
}

func (r *categoryRepository) GetModelsByBrandName(ctx context.Context, brandName string) ([]string, error) {
	var brandID gocql.UUID
	var modelList []string
	var model string
	brandQuery := "SELECT id FROM marketplace_keyspace.brands_by_name WHERE name = ?"
	if err := r.session.Query(brandQuery, brandName).WithContext(ctx).Scan(&brandID); err != nil {
		return nil, err
	}
	modelsQuery := "SELECT name FROM marketplace_keyspace.models WHERE brand_id = ? limit 5"
	iter := r.session.Query(modelsQuery, brandID).WithContext(ctx).Iter()
	for iter.Scan(&model) {
		modelList = append(modelList, model)
	}
	return modelList, nil
}

func (r *categoryRepository) GetBrandsByProductName(ctx context.Context, search string) ([]models.BrandWrap, error) {
	var brandList []models.BrandWrap
	var brandID gocql.UUID
	var productName string
	var categoryInfo models.CategoryInfo
	//var brandIDsList []gocql.UUID

	uniqueBrandIDs := make(map[gocql.UUID]bool)

	query := "SELECT brand_id, name FROM marketplace_keyspace.models WHERE name LIKE ? ALLOW FILTERING "
	iter := r.session.Query(query, search+"%").WithContext(ctx).Iter()
	for iter.Scan(&brandID, &productName) {

		if _, exists := uniqueBrandIDs[brandID]; exists {
			continue
		}
		uniqueBrandIDs[brandID] = true

		brandQuery := "SELECT name, subcategory_id FROM marketplace_keyspace.brands WHERE id = ?"
		brandIter := r.session.Query(brandQuery, brandID).Iter()
		var brandName string
		var subcategoryID gocql.UUID

		for brandIter.Scan(&brandName, &subcategoryID) {
			brand := models.BrandWrap{
				Name: brandName,
				ID:   brandID,
			}
			subcategoryQuery := "SELECT name, category_id FROM marketplace_keyspace.subcategories WHERE subcategory_id = ?"
			subcategoryIter := r.session.Query(subcategoryQuery, subcategoryID).WithContext(ctx).Iter()
			//Adding: subcategoryName, categoryID, categoryName
			for subcategoryIter.Scan(&categoryInfo.SubcategoryName, &categoryInfo.CategoryID) {
				categoryQuery := "SELECT name FROM marketplace_keyspace.category WHERE id = ?"
				if err := r.session.Query(categoryQuery, categoryInfo.CategoryID).WithContext(ctx).Scan(&categoryInfo.CategoryName); err != nil {
					return nil, err
				}
			}

			categoryInfo.SubcategoryID = subcategoryID
			brand.CategoryInfo = categoryInfo
			brandList = append(brandList, brand)
		}
	}
	return brandList, nil
}

func (r *categoryRepository) FieldsBySubcategory(ctx context.Context, subcategoryID gocql.UUID) (map[string]interface{}, error) {
	var filtersAndInputsJSON string

	query := "SELECT filtersandinputs FROM marketplace_keyspace.subcategories WHERE subcategory_id = ?"
	if err := r.session.Query(query, subcategoryID).WithContext(ctx).Scan(&filtersAndInputsJSON); err != nil {
		return nil, err
	}

	var filtersAndInputs map[string]interface{}
	if err := json.Unmarshal([]byte(filtersAndInputsJSON), &filtersAndInputs); err != nil {
		return nil, err
	}

	return filtersAndInputs, nil
}

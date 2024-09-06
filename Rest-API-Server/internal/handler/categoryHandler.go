package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/service"
	"marketplace_project/internal/utils"
	"net/http"
	"strings"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) AddCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	category.ID = gocql.TimeUUID()

	if err := h.service.AddCategory(&category); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *CategoryHandler) AddSubcategory(c *gin.Context) {
	var subcategoryParam models.SubcategoryParam

	if err := c.ShouldBindJSON(&subcategoryParam); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	subcategoryParam.Subcategory.ID = gocql.TimeUUID()

	if subcategoryParam.Subcategory.ParentID == nil && subcategoryParam.Subcategory.ParentName == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "ParentID or ParentName must be provided")
		return
	} else if subcategoryParam.Subcategory.ParentID == nil {
		if category, err := h.service.GetCategoryDataByName(context.Background(), subcategoryParam.Subcategory.ParentName); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Category with this name not exist")
			return
		} else {
			subcategoryParam.Subcategory.ParentID = &category.ID
			if err := h.service.AddSubcategory(&subcategoryParam); err != nil {
				utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	} else {
		if err := h.service.AddSubcategory(&subcategoryParam); err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func (h *CategoryHandler) ListMainCategories(c *gin.Context) {
	categories, err := h.service.ListMainCategories(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, categories)
}

func (h *CategoryHandler) ListAllCategories(c *gin.Context) {
	categories, err := h.service.ListAllCategories(context.Background())
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
	}
	utils.RespondWithJSON(c, http.StatusOK, categories)
}

type GroupedSubcategory struct {
	GroupID       gocql.UUID           `json:"id"`
	GroupName     string               `json:"groupName"`
	Subcategories []models.Subcategory `json:"subcategories"`
}

type CategoryResponse struct {
	ID                   gocql.UUID           `json:"id"`
	Name                 string               `json:"name"`
	Image                string               `json:"image"`
	GroupedSubcategories []GroupedSubcategory `json:"groupedSubcategories,omitempty"`
}

func (h *CategoryHandler) GroupSubcategoriesByCategory(c *gin.Context) {
	categories, err := h.service.ListAllCategories(context.Background())
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var response []CategoryResponse
	for _, category := range categories {
		groups := make(map[string][]models.Subcategory)
		for _, sub := range category.Subcategories {
			if sub.GroupName != "" {
				groups[sub.GroupName] = append(groups[sub.GroupName], sub)
			} else {
				groups["Uncategorized"] = append(groups["Uncategorized"], sub)
			}
		}

		var groupedSubcategories []GroupedSubcategory
		for groupName, subs := range groups {
			groupID := gocql.TimeUUID()
			groupedSubcategories = append(groupedSubcategories, GroupedSubcategory{
				GroupID:       groupID,
				GroupName:     groupName,
				Subcategories: subs,
			})
		}

		response = append(response, CategoryResponse{
			ID:                   category.ID,
			Name:                 category.Name,
			Image:                category.Image,
			GroupedSubcategories: groupedSubcategories,
		})
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *CategoryHandler) PopularCategories(c *gin.Context) {
	categories, err := h.service.PopularCategories(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, categories)
}

func (h *CategoryHandler) FieldsBySubcategory(c *gin.Context) {
	subcategoryID, err := gocql.ParseUUID(c.Query("subcategory_id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}
	fields, err := h.service.FieldsBySubcategory(context.Background(), subcategoryID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, fields)
}

func (h *CategoryHandler) BrandsBySubcategory(c *gin.Context) {
	subcategoryID, err := gocql.ParseUUID(c.Query("subcategory_id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid subcategory ID")
		return
	}
	brandsList, err := h.service.BrandsOfSubcategory(context.Background(), subcategoryID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, brandsList)
}

func (h *CategoryHandler) ModelsByBrands(c *gin.Context) {
	subcategoryID, err := gocql.ParseUUID(c.Query("brand_id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid brand ID")
		return
	}
	modelList, err := h.service.ModelsByBrands(context.Background(), subcategoryID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, modelList)
}

func (h *CategoryHandler) ParametersOfModels(c *gin.Context) {
	modelID, err := gocql.ParseUUID(c.Query("model_id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid model ID")
		return
	}
	parameters, err := h.service.ParametersOfModels(context.Background(), modelID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	type ParameterValue struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type Parameter struct {
		ParameterID   gocql.UUID       `json:"parameterID"`
		ParameterName string           `json:"parameterName"`
		DefaultOption string           `json:"defaultOption"`
		ParameterData []ParameterValue `json:"parameterData"`
	}

	var parametersSlice []Parameter

	for key, value := range parameters {
		var parameterData []ParameterValue
		for i, val := range value {
			parameterData = append(parameterData, ParameterValue{
				ID:   i + 1,
				Name: val,
			})
		}
		param := Parameter{
			ParameterID:   gocql.TimeUUID(),
			ParameterName: key,
			DefaultOption: strings.ToUpper(string(key[0])) + key[1:] + "s",
			ParameterData: parameterData,
		}
		parametersSlice = append(parametersSlice, param)
	}

	utils.RespondWithJSON(c, http.StatusOK, parametersSlice)
}

func (h *CategoryHandler) SearchByBrands(c *gin.Context) {
	search := c.Query("q")

	brandList, err := h.service.SearchByBrands(context.Background(), search)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid word")
		return
	}
	if brandList == nil {
		brandList, err = h.service.GetBrandsByProductName(context.Background(), search)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid word")
			return
		}
	}
	helpingWordsList, err := h.service.SearchByKeywords(context.Background(), search)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid word")
		return
	}

	if helpingWordsList == nil && len(brandList) > 0 {
		for _, brand := range brandList {
			helpingWordsList, err = h.service.GetModelsByBrandName(context.Background(), brand.Name)
			if err != nil {
				utils.RespondWithError(c, http.StatusBadRequest, "Invalid word")
				return
			}
		}
	}

	type Keyword struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var keywords []Keyword

	for i, keyword := range helpingWordsList {
		keywords = append(keywords, Keyword{
			ID:   i + 1,
			Name: keyword,
		})
	}

	repository := map[string]interface{}{
		"brands":   brandList,
		"keywords": keywords,
	}

	utils.RespondWithJSON(c, http.StatusOK, repository)
}

func (h *CategoryHandler) InsertSubcategoriesToGroup(c *gin.Context) {
	var group *models.SubcategoryGroups
	if err := c.ShouldBindJSON(&group); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.service.InsertSubcategoriesToGroup(context.Background(), group); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, "Group is inserted")
}

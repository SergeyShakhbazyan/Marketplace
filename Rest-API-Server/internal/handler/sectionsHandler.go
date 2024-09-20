package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/service"
	"marketplace_project/internal/utils"
	"net/http"
)

type SectionsHandler struct {
	service *service.SectionsService
}

func NewSectionsHandler(service *service.SectionsService) *SectionsHandler {
	return &SectionsHandler{service: service}
}

func (h *SectionsHandler) GetMainCategoriesSection(c *gin.Context) {
	categoriesSection, err := h.service.MainCategoriesSection(context.Background())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, categoriesSection)
}

func (h *SectionsHandler) GetMainProductsSections(c *gin.Context) {
	productsSection, err := h.service.MainProductsSections(context.Background())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, productsSection)
}

func (h *SectionsHandler) Section(c *gin.Context) {
	pageName := c.Query("pageName")
	if pageName == "home" {
		categoriesSection, err := h.service.MainCategoriesSection(context.Background())
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		productsSection, err := h.service.MainProductsSections(context.Background())
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		response := gin.H{
			"categoriesSection": categoriesSection,
			"productSections":   productsSection,
		}
		utils.RespondWithJSON(c, http.StatusOK, response)
	}

	//if pageName == "profileAds" {
	//	profileAds, err = h.service.GetUserProducts(context.Background(), ownerId)
	//}

}

func (h *SectionsHandler) GetProfileInfo(c *gin.Context) {
	userID, err := gocql.ParseUUID(c.Query("userID"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, products, err := h.service.GetProfileInfo(context.Background(), userID)

	type UserData struct {
		UserID    gocql.UUID                  `json:"userID"`
		Avatar    string                      `json:"avatar"`
		FirstName string                      `json:"firstName"`
		LastName  string                      `json:"lastName"`
		Products  []models.ProductWrapContent `json:"products"`
	}

	response := UserData{
		UserID:    user.UserID,
		Avatar:    user.Avatar,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Products:  products,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

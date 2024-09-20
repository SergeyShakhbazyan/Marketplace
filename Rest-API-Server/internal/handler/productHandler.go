package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/service"
	"marketplace_project/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func extractKeywords(description string /*, tags []string8*/) []string {
	description = strings.ToLower(description)

	splitFunc := func(r rune) bool {
		return r == ' ' || r == '_' || r == '-' || r == '.' || r == ','
	}

	words := strings.FieldsFunc(description, splitFunc)

	keywordSet := make(map[string]struct{})
	for _, word := range words {
		keywordSet[word] = struct{}{}
	}

	//for _, tag := range tags {
	//	tag = strings.ToLower(tag)
	//	keywordSet[tag] = struct{}{}
	//}

	var keywords []string
	for word := range keywordSet {
		keywords = append(keywords, word)
	}

	return keywords
}

type ProductRequest struct {
	Product models.Product      `json:"product"`
	Filters []map[string]string `json:"filters"`
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	req.Product.ProductID = gocql.TimeUUID()
	req.Product.Keywords = extractKeywords(req.Product.Title /*, product.Tags*/)
	req.Product.CreatedAt = time.Now()

	if err := h.service.AddProduct(&req.Product, &req.Filters); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := gocql.ParseUUID(c.Query("productID"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = h.service.DeleteProduct(productID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, "Product deleted")
}

func (h *ProductHandler) ProductInfo(c *gin.Context) {
	productID, err := gocql.ParseUUID(c.Query("productID"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}
	productInfo, filters, err := h.service.ProductInfoByID(productID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	response := map[string]interface{}{
		"productInfo": productInfo,
		"filters":     filters,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *ProductHandler) ProductsByCategoryBeta(c *gin.Context) {
	categoryID, err := gocql.ParseUUID(c.Query("category_id"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	lastProductIDStr := c.Query("lastProductID")
	var lastProductID gocql.UUID
	if lastProductIDStr != "" {
		lastProductID, err = gocql.ParseUUID(lastProductIDStr)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid last product Time")
			return
		}
	}

	products, newPagingStateStr, err := h.service.ProductsWrapsByCategory(categoryID, lastProductID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagingState": newPagingStateStr,
		"products":    products,
	})
}

func (h *ProductHandler) ProductsByOwnerID(c *gin.Context) {
	userID, err := gocql.ParseUUID(c.Query("userID"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	products, err := h.service.GetProductsByOwnerID(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, products)
}

func (h *ProductHandler) Products(c *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, products)
}

func (h *ProductHandler) SearchEngine(c *gin.Context) {
	searchQuery := c.Query("search")
	searchQuery = strings.ToLower(searchQuery)
	searchQuery = strings.ReplaceAll(searchQuery, "+", " ")
	searchProduct, err := h.service.SearchProducts(searchQuery)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, searchProduct)
}

//func (h *ProductHandler) FindProductsByFilters(c *gin.Context) {
//	queryParams := c.Request.URL.Query()
//	var categoryID gocql.UUID
//	var subcategoryID gocql.UUID
//	filters := make(map[string]string)
//	for key, values := range queryParams {
//		if len(values) > 0 {
//			if key == "category" {
//				categoryID, _ = gocql.ParseUUID(values[0])
//			} else if key == "subcategory" {
//				subcategoryID, _ = gocql.ParseUUID(values[0])
//			} else {
//				filters[key] = values[0]
//			}
//		}
//	}
//
//	productIDs := make(map[gocql.UUID]bool)
//	var intersection []gocql.UUID
//
//	for filterName, filterValue := range filters {
//
//		var filter = models.Filter{
//			FilterName:  filterName,
//			FilterValue: filterValue,
//		}
//
//		ids, err := h.service.FindProductsByFilters(categoryID, subcategoryID, filter)
//		if err != nil {
//			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
//			return
//		}
//
//		if len(productIDs) == 0 {
//			for _, id := range ids {
//				productIDs[id] = true
//			}
//		} else {
//			tempIDs := make(map[gocql.UUID]bool)
//			for _, id := range ids {
//				if _, exists := productIDs[id]; exists {
//					tempIDs[id] = true
//				}
//			}
//			productIDs = tempIDs
//		}
//	}
//	for id := range productIDs {
//		intersection = append(intersection, id)
//	}
//}

func (h *ProductHandler) FindProductsByFilters(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var categoryID, subcategoryID gocql.UUID
	filters := make(map[string]string)
	limit := 0

	for key, values := range queryParams {
		if len(values) > 0 {
			switch key {
			case "category":
				var err error
				categoryID, err = gocql.ParseUUID(values[0])
				if err != nil {
					utils.RespondWithError(c, http.StatusBadRequest, "Invalid category UUID")
					return
				}
			case "subcategory":
				var err error
				subcategoryID, err = gocql.ParseUUID(values[0])
				if err != nil {
					utils.RespondWithError(c, http.StatusBadRequest, "Invalid subcategory UUID")
					return
				}

			case "limit":
				var err error
				limit, err = strconv.Atoi(values[0])
				if err != nil || limit < 0 {
					utils.RespondWithError(c, http.StatusBadRequest, "Invalid limit value")
					return
				}
			default:
				filters[key] = values[0]
			}
		}
	}

	if len(filters) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "No filters provided")
		return
	}

	// Map to store intersecting product IDs
	productIDs := make(map[gocql.UUID]bool)

	// Perform the filtering and intersection logic
	for filterName, filterValue := range filters {
		filter := models.Filter{
			Name:  filterName,
			Value: filterValue,
		}

		ids, err := h.service.FindProductsByFilters(categoryID, subcategoryID, filter, limit)
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		if len(productIDs) == 0 {
			// Initialize the map with the first filter's result
			for _, id := range ids {
				productIDs[id] = true
			}
		} else {
			// Perform intersection with existing productIDs
			tempIDs := make(map[gocql.UUID]bool)
			for _, id := range ids {
				if productIDs[id] {
					tempIDs[id] = true
				}
			}
			productIDs = tempIDs
		}

		// If at any point the intersection is empty, break early
		if len(productIDs) == 0 {
			break
		}
	}

	// Convert the map keys to a slice
	var intersection []gocql.UUID
	for id := range productIDs {
		intersection = append(intersection, id)
	}

	var products []models.ProductWrapContent

	for _, id := range intersection {
		product, err := h.service.FindProductsByIDs(context.Background(), id)
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		products = append(products, *product)
	}

	utils.RespondWithJSON(c, http.StatusOK, products)
}

package app

import (
	"marketplace_project/internal/handler"
)

func (a *App) setRoutersForUser(userHandler *handler.UserHandler) {
	a.Router.POST("/register", userHandler.Register)
	a.Router.POST("/signIn", userHandler.SignIn)
	a.Router.POST("/token", userHandler.RefreshToken)
	a.Router.GET("/profileData", userHandler.UserDataByID)
}

func (a *App) setRoutersForCategory(categoryHandler *handler.CategoryHandler) {
	a.Router.POST("/addCategory", categoryHandler.AddCategory)
	a.Router.POST("/addGroup", categoryHandler.InsertSubcategoriesToGroup)
	a.Router.GET("/mainCategories", categoryHandler.ListMainCategories)
	a.Router.GET("/categories", categoryHandler.ListAllCategories)
	a.Router.GET("/catalog", categoryHandler.GroupSubcategoriesByCategory)
	a.Router.POST("/addSubcategory", categoryHandler.AddSubcategory)
	a.Router.GET("/popularCategories", categoryHandler.PopularCategories)
	a.Router.GET("/subcategoryFields", categoryHandler.FieldsBySubcategory)
	a.Router.GET("/brands", categoryHandler.BrandsBySubcategory)
	a.Router.GET("/models", categoryHandler.ModelsByBrands)
	a.Router.GET("/parameters", categoryHandler.ParametersOfModels)
	a.Router.GET("/search", categoryHandler.SearchByBrands)
}

func (a *App) setRoutersForProduct(productHandler *handler.ProductHandler) {
	a.Router.POST("/addProduct", productHandler.AddProduct)
	a.Router.DELETE("/deleteProduct", productHandler.DeleteProduct)
	a.Router.POST("/recommendedProducts", productHandler.FindProductsByFilters)
	a.Router.GET("/products", productHandler.Products)
	a.Router.GET("/productsByCategory", productHandler.ProductsByCategoryBeta)
	a.Router.GET("/searchProduct", productHandler.SearchEngine)
	a.Router.GET("/findProduct", productHandler.FindProductsByFilters)
	a.Router.GET("/product", productHandler.ProductInfo)
}

func (a *App) setRoutersForSections(sectionHandler *handler.SectionsHandler) {
	a.Router.GET("/getPageSections", sectionHandler.Section)
	a.Router.GET("/user", sectionHandler.GetProfileInfo)
}

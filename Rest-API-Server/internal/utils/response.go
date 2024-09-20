package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrCategoryExists = errors.New("category already exists")
	ErrNotFound       = errors.New("record not found")
)

func RespondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

func RespondWithError(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(statusCode, gin.H{"error": errMsg})
}

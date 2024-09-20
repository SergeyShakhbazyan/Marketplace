package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/service"
	"marketplace_project/internal/utils"
	"net/http"
	"time"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(serviceUser *service.UserService) *UserHandler {
	return &UserHandler{service: serviceUser}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.UserID = gocql.TimeUUID()
	user.CreatedAt = time.Now()
	user.Avatar = "https://firebasestorage.googleapis.com/v0/b/marketplace-dee62.appspot.com/o/avatars%2Favatar-default.svg?alt=media&token=ee6f1132-fa12-4be1-ad5c-338487892508"
	if err := h.service.Register(&user); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(c, http.StatusCreated, user)
}

type option struct {
	Exp int `json:"exp"`
}

type Token struct {
	Token   string `json:"token"`
	Options option `json:"options"`
}

type tokens struct {
	AccessToken  Token `json:"accessToken"`
	RefreshToken Token `json:"refreshToken"`
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.service.SignIn(credentials.Email, credentials.Password)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid email or password")
		return
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.UserID, user.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	userData := models.UserWrapContent{
		UserID:      user.UserID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AccountType: user.AccountType,
		Avatar:      user.Avatar,
	}

	accessTokenOption := Token{
		Token:   accessToken,
		Options: option{Exp: 10000},
	}

	refreshTokenOption := Token{
		Token:   refreshToken,
		Options: option{Exp: 604800000},
	}

	tokens := tokens{
		AccessToken:  accessTokenOption,
		RefreshToken: refreshTokenOption,
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
		"user":   userData,
	})
}

type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var tokenRequest TokenRequest
	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, email, err := utils.ExtractUserIdAndEmailFromContext(tokenRequest.Token)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	accessToken, refreshToken, err := utils.GenerateToken(userID, email)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate new token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
func (h *UserHandler) UserDataByID(c *gin.Context) {
	userID, err := gocql.ParseUUID(c.Query("userID"))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	userData, err := h.service.GetUserByID(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, userData)
}

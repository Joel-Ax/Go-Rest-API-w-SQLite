package routes

import (
	"net/http"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.BindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse request data": err.Error()})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not save user": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"User was created:": user})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse request data": err.Error()})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Could not validate credentials": err.Error()})
		return
	}
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not validate credentials": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "token": token})
}

func logout(context *gin.Context) {

	context.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

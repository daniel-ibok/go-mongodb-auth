package controllers

import (
	"go-mongodb-auth/jwt"
	"go-mongodb-auth/models"
	"go-mongodb-auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginController(c *gin.Context) {
	var user models.LoginUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
		return
	}

	// check if user already exists
	if exists := models.CheckUser(user.Email); !exists {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": "user does not exists",
		})
		return
	}

	// retrieve user
	authUser, _ := models.GetUserByEmail(user.Email)

	// compare password
	if err := utils.CheckPassword(user.Password, authUser.Password); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": "Password is incorrect",
		})
		return
	}

	// create and store token in cookie
	token, _, err := jwt.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	utils.GenerateAndStoreCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"accessToken": token,
			"user": gin.H{
				"username": authUser.Username,
				"email":    authUser.Email,
			},
		},
	})

}

func RegisterController(c *gin.Context) {
	var user models.RegisterUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
		return
	}

	// check if user already exists
	if exists := models.CheckUser(user.Email); exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "user does not exists",
		})
		return
	}

	// hash password
	user.Password, _ = utils.HashPassword(user.Password)

	// register new user
	err := models.NewUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// create and store token in cookie
	token, _, err := jwt.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	utils.GenerateAndStoreCookie(c, token)

	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"data": gin.H{
			"accessToken": token,
			"user": gin.H{
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

func UserAuthController(c *gin.Context) {
	data, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Cannot retrieve user",
		})
		return
	}

	user := data.(*models.User)
	c.JSON(http.StatusOK, user)
}

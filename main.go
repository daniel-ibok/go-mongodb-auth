package main

import (
	"go-mongodb-auth/controllers"
	"go-mongodb-auth/database"
	"go-mongodb-auth/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	// initialize database connection
	database.NewDBInstance()

	// setup routes
	router := gin.Default()
	router.POST("/login", controllers.LoginController)
	router.POST("/register", controllers.RegisterController)
	router.GET("/dashboard", middleware.AuthMiddleware(), controllers.UserAuthController)
	router.Run(":8005")
}

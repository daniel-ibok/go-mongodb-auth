package main

import (
	"go-mongodb-auth/controllers"
	"go-mongodb-auth/database"
	"go-mongodb-auth/middleware"
	"net/http"
	"time"

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

	server := &http.Server{
		Addr:         ":8005",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

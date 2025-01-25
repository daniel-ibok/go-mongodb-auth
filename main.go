package main

import (
	"context"
	"fmt"
	"go-mongodb-auth/controllers"
	"go-mongodb-auth/database"
	"go-mongodb-auth/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// initialize database connection
	database.NewDBInstance()
	done := make(chan bool)

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

	// shutdown server
	go ShutdownServer(server, done)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		panic(err)
	}

	// wait for server to gracefully exit
	<-done
}

func ShutdownServer(server *http.Server, done chan bool) {
	// create signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown error", err)
		return
	}

	// close database connection
	if err := database.Close(); err != nil {
		log.Println("Database shutdown error", err)
	}

	<-ctx.Done()
	log.Println("Server exiting...")
	done <- true
}

package main

import (
	"be-jalan/config"
	"be-jalan/handlers"
	"be-jalan/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{})
	r := gin.Default()

	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/verify-email", handlers.VerifyEmail)
		authGroup.POST("/verify-phone", handlers.VerifyPhone)
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/forgot-password", handlers.ForgotPassword)
		authGroup.POST("/verify-otp", handlers.VerifyOTP)
		authGroup.POST("/reset-password", handlers.ResetPassword)
		authGroup.POST("/google", handlers.GoogleAuth)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

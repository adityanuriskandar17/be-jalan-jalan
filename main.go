package main

import (
	"be-jalan/config"
	"be-jalan/database"
	"be-jalan/handlers"
	"be-jalan/middleware"
	"be-jalan/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Destination{},
		&models.DestinationImage{},
		&models.Ticket{},
		&models.Review{},
		&models.Favorite{},
		&models.Notification{},
		&models.Order{},
		&models.OrderItem{},
		&models.PromoCode{},
	)

	// Seed Initial Data (Admin Users & Travel Data)
	database.SeedUsers(config.DB)
	database.SeedTravelData(config.DB)

	r := gin.Default()
	r.Static("/uploads", "./uploads")

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

	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.AuthMiddleware())
	{
		adminGroup.POST("/upload", handlers.UploadImage)
		adminGroup.GET("/dashboard", handlers.GetDashboardStats)

		// Destinations
		adminGroup.GET("/destinations/stats", handlers.GetDestinationStats)
		adminGroup.GET("/destinations", handlers.GetDestinations)
		adminGroup.POST("/destinations", handlers.CreateDestination)
		adminGroup.PUT("/destinations/:id", handlers.UpdateDestination)
		adminGroup.PATCH("/destinations/:id/status", handlers.ToggleDestinationStatus)
		adminGroup.DELETE("/destinations/:id", handlers.DeleteDestination)

		// Categories
		adminGroup.GET("/categories", handlers.GetCategories)
		adminGroup.POST("/categories", handlers.CreateCategory)
		adminGroup.PUT("/categories/:id", handlers.UpdateCategory)
		adminGroup.DELETE("/categories/:id", handlers.DeleteCategory)

		// Tickets
		adminGroup.GET("/tickets/stats", handlers.GetTicketStats)
		adminGroup.GET("/tickets", handlers.GetTickets)
		adminGroup.POST("/tickets", handlers.CreateTicket)
		adminGroup.PUT("/tickets/:id", handlers.UpdateTicket)
		adminGroup.DELETE("/tickets/:id", handlers.DeleteTicket)

		// Bookings
		adminGroup.GET("/bookings/stats", handlers.GetBookingStats)
		adminGroup.GET("/bookings", handlers.GetBookings)
		adminGroup.GET("/bookings/:id", handlers.GetBookingDetail)
		adminGroup.GET("/bookings/:id/eticket", handlers.GetBookingETicket)
		adminGroup.POST("/bookings/scan", handlers.ScanTicket) // New: Scan QR
		adminGroup.POST("/bookings", handlers.CreateManualBooking)
		adminGroup.PUT("/bookings/:id/status", handlers.UpdateBookingStatus)
		adminGroup.DELETE("/bookings/:id", handlers.DeleteBooking)

		// Users
		adminGroup.GET("/users/stats", handlers.GetUserStats)
		adminGroup.GET("/users", handlers.GetUsers)
		adminGroup.PATCH("/users/:id/status", handlers.ToggleUserStatus)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

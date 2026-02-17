package database

import (
	"be-jalan/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	// Define roles
	const (
		RoleUser       = "user"
		RoleSuperAdmin = "super_admin"
		RoleAdmin      = "admin"
		RoleCS         = "cs"
	)

	// List of users to seed
	users := []models.User{
		{
			FullName:        "Super Admin",
			Email:           "admin@jalanjalanyuk.com",
			Password:        "admin123", // Will be hashed
			PhoneNumber:     "081234567890",
			Role:            RoleSuperAdmin,
			IsEmailVerified: true,
			IsPhoneVerified: true,
		},
		{
			FullName:        "Customer Service",
			Email:           "cs@jalanjalanyuk.com",
			Password:        "cs123",
			PhoneNumber:     "081234567891",
			Role:            RoleCS,
			IsEmailVerified: true,
			IsPhoneVerified: true,
		},
		{
			FullName:        "Admin Operator",
			Email:           "operator@jalanjalanyuk.com",
			Password:        "operator123",
			PhoneNumber:     "081234567892",
			Role:            RoleAdmin,
			IsEmailVerified: true,
			IsPhoneVerified: true,
		},
	}

	for _, user := range users {
		// Check if user exists
		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			fmt.Printf("User %s already exists, skipping...\n", user.Email)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		user.Password = string(hashedPassword)

		// Create user
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
		} else {
			fmt.Printf("✅ Created user: %s (%s)\n", user.Email, user.Role)
		}
	}
}

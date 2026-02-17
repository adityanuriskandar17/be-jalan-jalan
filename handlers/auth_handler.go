package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"be-jalan/utils"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Helper function to check password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Helper function to generate 6-digit OTP
func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return padLeft(n.String(), "0", 6), nil
}

// padLeft ensures the OTP is always 6 digits
func padLeft(str, pad string, length int) string {
	for len(str) < length {
		str = pad + str
	}
	return str
}

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if result := config.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Check if phone number already exists
	if result := config.DB.Where("phone_number = ?", req.PhoneNumber).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Generate OTPs
	emailOTP, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Email OTP"})
		return
	}

	phoneOTP, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Phone OTP"})
		return
	}

	expiryTime := time.Now().Add(24 * time.Hour) // 24 hours for registration verification? or 15 mins? sticking to 15 mins for consistency or 24h as per common practice for links, but OTPs usually short lived. Let's do 15 mins.
	expiryTime = time.Now().Add(15 * time.Minute)

	newUser := models.User{
		FullName:    req.FullName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,

		EmailVerificationToken:  emailOTP,
		EmailVerificationExpiry: &expiryTime,
		IsEmailVerified:         false,

		PhoneVerificationToken:  phoneOTP,
		PhoneVerificationExpiry: &expiryTime,
		IsPhoneVerified:         false,
	}

	if result := config.DB.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// TODO: Send Email OTP
	// TODO: Send SMS OTP

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil. Silakan verifikasi email dan nomor telepon Anda.",
		"data": gin.H{
			"id":    newUser.ID,
			"email": newUser.Email,
		},
		"debug_email_otp": emailOTP, // Remove in production
		"debug_phone_otp": phoneOTP, // Remove in production
	})
}

// VerifyEmail verifies the Email OTP
func VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.IsEmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already verified"})
		return
	}

	if user.EmailVerificationToken != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if user.EmailVerificationExpiry != nil && time.Now().After(*user.EmailVerificationExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	user.IsEmailVerified = true
	user.EmailVerificationToken = ""
	user.EmailVerificationExpiry = nil

	if result := config.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verifcation status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// VerifyPhone verifies the Phone OTP
func VerifyPhone(c *gin.Context) {
	var req models.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("phone_number = ?", req.PhoneNumber).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.IsPhoneVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Phone number already verified"})
		return
	}

	if user.PhoneVerificationToken != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if user.PhoneVerificationExpiry != nil && time.Now().After(*user.PhoneVerificationExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	user.IsPhoneVerified = true
	user.PhoneVerificationToken = ""
	user.PhoneVerificationExpiry = nil

	if result := config.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Phone number verified successfully",
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !checkPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data": gin.H{
			"token": token,
			"user": models.User{
				ID:       user.ID,
				FullName: user.FullName,
				Email:    user.Email,
				Role:     user.Role,
			},
		},
	})
}

// ForgotPassword handles password reset request
func ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		// Verify if user exists. If not found, return generic success message for security but log it
		// For development, returning 404 is fine or error message
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	otp, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// Set OTP and expiration (15 minutes)
	expiryTime := time.Now().Add(15 * time.Minute)
	user.ForgotPasswordToken = otp
	user.ForgotPasswordTokenExpiry = &expiryTime

	if result := config.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}

	// TODO: Send OTP via email (using SMTP or transactional email service)
	// For now, return it in response for testing purposes (remove in production!)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Kode OTP untuk reset password telah dikirim ke email Anda.",
		"debug_otp": otp, // Remove this in production!!
	})
}

// VerifyOTP verifies the OTP provided by the user
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ForgotPasswordToken != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if user.ForgotPasswordTokenExpiry != nil && time.Now().After(*user.ForgotPasswordTokenExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP valid",
	})
}

// ResetPassword handles the actual password reset process
func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Re-verify OTP for security (stateless check)
	if user.ForgotPasswordToken != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	if user.ForgotPasswordTokenExpiry != nil && time.Now().After(*user.ForgotPasswordTokenExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	// Hash new password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Update password and clear OTP
	user.Password = hashedPassword
	user.ForgotPasswordToken = ""
	user.ForgotPasswordTokenExpiry = nil

	if result := config.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// GoogleAuth handles login via Google
func GoogleAuth(c *gin.Context) {
	var req models.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Verify Google token and create/login user

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Google berhasil",
		"data": gin.H{
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mocktoken...",
			"user": models.User{
				ID:       123, // Use int/uint
				FullName: "Google User",
				Email:    "google@example.com",
				Role:     "user",
			},
		},
	})
}

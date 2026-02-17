package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"be-jalan/services/notification"
	"be-jalan/utils"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
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

// Helper function to hash OTP
func hashOTP(otp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(otp), 14)
	return string(bytes), err
}

// Helper function to check OTP hash
func checkOTPHash(otp, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(otp))
	return err == nil
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
	hashEmailOTP, err := hashOTP(emailOTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process Email OTP"})
		return
	}

	phoneOTP, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Phone OTP"})
		return
	}
	hashPhoneOTP, err := hashOTP(phoneOTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process Phone OTP"})
		return
	}

	expiryTime := time.Now().Add(15 * time.Minute)

	newUser := models.User{
		FullName:    req.FullName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,

		EmailVerificationToken:  hashEmailOTP,
		EmailVerificationExpiry: &expiryTime,
		IsEmailVerified:         false,

		PhoneVerificationToken:  hashPhoneOTP,
		PhoneVerificationExpiry: &expiryTime,
		IsPhoneVerified:         false,
	}

	if result := config.DB.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send Email OTP
	emailConfig := notification.GetEmailConfig()
	fmt.Printf("📧 Sending Email using SMTP Host: %s\n", emailConfig.SMTPHost)
	if err := notification.SendOTPEmail(emailConfig, req.Email, emailOTP, "Verifikasi Email"); err != nil {
		// Log error but don't fail registration
		fmt.Printf("⚠️ Failed to send email OTP: %v\n", err)
	}

	// Send WhatsApp OTP
	tokenWA := os.Getenv("WA_API_TOKEN")
	secretWA := os.Getenv("WA_API_SECRET")
	cfgWA := notification.WhatsAppConfig{
		APIToken:  tokenWA,
		APISecret: secretWA,
	}

	if err := notification.SendWhatsAppOTP(cfgWA, req.PhoneNumber, phoneOTP); err != nil {
		fmt.Printf("❌ Gagal kirim WhatsApp: %v\n", err)
	}

	// Optional: Send SMS OTP as backup
	if err := notification.SendSMSOTP(notification.GetSMSConfig(), req.PhoneNumber, phoneOTP); err != nil {
		// Log error but don't fail registration
		fmt.Printf("⚠️ Failed to send SMS OTP: %v\n", err)
	}

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

	if !checkOTPHash(req.OTP, user.EmailVerificationToken) {
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

	if !checkOTPHash(req.OTP, user.PhoneVerificationToken) {
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
	hashOTP, err := hashOTP(otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process OTP"})
		return
	}

	// Set OTP and expiration (15 minutes)
	expiryTime := time.Now().Add(15 * time.Minute)
	user.ForgotPasswordToken = hashOTP
	user.ForgotPasswordTokenExpiry = &expiryTime

	if result := config.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}

	// Send OTP via Email or WhatsApp/SMS based on identifier
	if strings.Contains(req.Email, "@") {
		// Send via Email
		if err := notification.SendOTPEmail(notification.GetEmailConfig(), user.Email, otp, "Reset Password"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email OTP"})
			return
		}
	} else {
		// Send via WhatsApp
		if err := notification.SendWhatsAppOTP(notification.GetWhatsAppConfig(), user.PhoneNumber, otp); err != nil {
			// Try SMS as fallback
			if err := notification.SendSMSOTP(notification.GetSMSConfig(), user.PhoneNumber, otp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Kode OTP untuk reset password telah dikirim.",
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

	if !checkOTPHash(req.OTP, user.ForgotPasswordToken) {
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
	if !checkOTPHash(req.OTP, user.ForgotPasswordToken) {
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

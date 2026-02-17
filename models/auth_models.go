package models

import (
	"time"
)

// RegisterRequest represents the JSON body for user registration
type RegisterRequest struct {
	FullName    string `json:"fullName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required,numeric"`
	Password    string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the JSON body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ForgotPasswordRequest represents the JSON body for password reset request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyOTPRequest represents the JSON body for verifying OTP
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// ResetPasswordRequest represents the JSON body for resetting password
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	OTP             string `json:"otp" binding:"required,len=6"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

// GoogleAuthRequest represents the JSON body for Google OAuth login
type GoogleAuthRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

// VerifyEmailRequest represents the JSON body for verifying email
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// VerifyPhoneRequest represents the JSON body for verifying phone number
type VerifyPhoneRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,numeric"`
	OTP         string `json:"otp" binding:"required,len=6"`
}

// User is the main user model
type User struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	FullName    string `json:"fullName"`
	Email       string `gorm:"uniqueIndex" json:"email"`
	PhoneNumber string `gorm:"uniqueIndex" json:"phoneNumber"`
	Password    string `json:"-"` // Don't expose password in JSON
	Role        string `json:"role" gorm:"default:'user'"`

	// Forgot Password
	ForgotPasswordToken       string     `json:"-"`
	ForgotPasswordTokenExpiry *time.Time `json:"-"`

	// Email Verification
	IsEmailVerified         bool       `gorm:"default:false" json:"isEmailVerified"`
	EmailVerificationToken  string     `json:"-"`
	EmailVerificationExpiry *time.Time `json:"-"`

	// Phone Verification
	IsPhoneVerified         bool       `gorm:"default:false" json:"isPhoneVerified"`
	PhoneVerificationToken  string     `json:"-"`
	PhoneVerificationExpiry *time.Time `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

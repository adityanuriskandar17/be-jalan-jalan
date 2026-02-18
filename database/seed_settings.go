package database

import (
	"be-jalan/models"

	"gorm.io/gorm"
)

// SeedSettings populates default system configurations
func SeedSettings(db *gorm.DB) {
	// Defaults
	defaults := []models.Setting{
		// Payment Methods (bool)
		{Group: "payment", Key: "payment_qris_enabled", Value: "true", Type: "bool"},
		{Group: "payment", Key: "payment_ovo_enabled", Value: "true", Type: "bool"},
		{Group: "payment", Key: "payment_gopay_enabled", Value: "true", Type: "bool"},
		{Group: "payment", Key: "payment_dana_enabled", Value: "true", Type: "bool"},
		{Group: "payment", Key: "payment_bca_enabled", Value: "true", Type: "bool"},
		{Group: "payment", Key: "payment_mandiri_enabled", Value: "true", Type: "bool"},

		// Email Settings (bool)
		{Group: "email", Key: "email_notification_enabled", Value: "true", Type: "bool"},
		{Group: "email", Key: "email_booking_confirm_enabled", Value: "true", Type: "bool"},
		{Group: "email", Key: "email_payment_reminder_enabled", Value: "true", Type: "bool"},
		{Group: "email", Key: "email_promo_enabled", Value: "false", Type: "bool"},

		// System (Security/Access) (bool)
		{Group: "system", Key: "system_maintenance_mode", Value: "false", Type: "bool"},
		{Group: "system", Key: "system_user_registration_enabled", Value: "true", Type: "bool"},
		{Group: "system", Key: "system_email_verification_required", Value: "true", Type: "bool"},

		// General (string)
		{Group: "general", Key: "app_name", Value: "JalanJalanYuk", Type: "string"},
		{Group: "general", Key: "support_email", Value: "support@jalanjalanyuk.com", Type: "string"},
		{Group: "general", Key: "support_phone", Value: "+62 812-3456-7890", Type: "string"},
		{Group: "general", Key: "timezone", Value: "Asia/Jakarta (WIB)", Type: "string"},
	}

	for _, s := range defaults {
		var existing models.Setting
		// Check if exists, if not create
		if err := db.Where("key = ?", s.Key).First(&existing).Error; err != nil {
			db.Create(&s)
		}
	}
}

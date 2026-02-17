package notification

import (
	"os"
	"strconv"
)

// GetEmailConfig returns EmailConfig from environment variables
func GetEmailConfig() EmailConfig {
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
		FromName:     os.Getenv("SMTP_FROM_NAME"),
	}
}

// GetWhatsAppConfig returns WhatsAppConfig from environment variables
func GetWhatsAppConfig() WhatsAppConfig {
	return WhatsAppConfig{
		APIToken:  os.Getenv("WA_API_TOKEN"),
		APISecret: os.Getenv("WA_API_SECRET"),
	}
}

// GetSMSConfig returns SMSConfig from environment variables
func GetSMSConfig() SMSConfig {
	return SMSConfig{
		BaseURL:  os.Getenv("SMS_BASE_URL"),
		Username: os.Getenv("SMS_USERNAME"),
		Password: os.Getenv("SMS_PASSWORD"),
	}
}

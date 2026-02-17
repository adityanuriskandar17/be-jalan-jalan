package notification

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// SendEmail sends an email using SMTP
func SendEmail(cfg EmailConfig, to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", cfg.FromName, cfg.FromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("✅ Email sent to %s\n", to)
	return nil
}

// SendOTPEmail sends an OTP verification email
func SendOTPEmail(cfg EmailConfig, to, otp, purpose string) error {
	subject := fmt.Sprintf("Kode OTP - %s", purpose)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 30px; border-radius: 5px; margin-top: 20px; }
				.otp-code { font-size: 32px; font-weight: bold; color: #4CAF50; text-align: center; padding: 20px; background-color: white; border-radius: 5px; margin: 20px 0; letter-spacing: 5px; }
				.footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
				.warning { color: #f44336; font-weight: bold; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>ASITA - Jalan Jalan</h1>
				</div>
				<div class="content">
					<h2>%s</h2>
					<p>Kode OTP Anda adalah:</p>
					<div class="otp-code">%s</div>
					<p>Kode ini berlaku selama <strong>15 menit</strong>.</p>
					<p class="warning">⚠️ Jangan bagikan kode ini kepada siapapun!</p>
					<p>Jika Anda tidak meminta kode ini, abaikan email ini.</p>
				</div>
				<div class="footer">
					<p>&copy; 2026 ASITA - Jalan Jalan. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, purpose, otp)

	return SendEmail(cfg, to, subject, htmlBody)
}

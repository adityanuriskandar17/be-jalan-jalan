package notification

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type SMSConfig struct {
	BaseURL  string
	Username string
	Password string
}

// SendSMSOTP sends an OTP via SMS using GoSMSGateway Masking API
func SendSMSOTP(cfg SMSConfig, mobile, otp string) error {
	message := fmt.Sprintf("Kode OTP Anda: %s. Berlaku 15 menit. Jangan bagikan kode ini!", otp)
	return SendSMS(cfg, mobile, message, 0)
}

// SendSMS sends a custom SMS message using GoSMSGateway Masking API
// lineType: 0 = reguler, 1 = prioritas
func SendSMS(cfg SMSConfig, mobile, message string, lineType int) error {
	// Prepare form data
	data := url.Values{}
	data.Set("username", cfg.Username)
	data.Set("mobile", mobile)
	data.Set("message", message)
	data.Set("type", strconv.Itoa(lineType))

	// Build URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", cfg.BaseURL, data.Encode())

	// Create request with basic auth
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set basic authentication
	req.SetBasicAuth(cfg.Username, cfg.Password)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("✅ SMS sent to %s. Response: %s\n", mobile, string(bodyBytes))
	return nil
}

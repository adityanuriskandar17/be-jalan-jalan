package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WhatsAppConfig struct {
	APIToken  string
	APISecret string
}

// SendWhatsAppOTP sends an OTP via WhatsApp using specific API
func SendWhatsAppOTP(cfg WhatsAppConfig, mobile, otp string) error {
	// Use URL from environment or fallback
	apiURL := os.Getenv("WA_URL_INVOICE")
	if apiURL == "" {
		apiURL = "https://apiwa.asitatech.co.id/api/sendMessageText"
	}

	// Prepare message
	message := fmt.Sprintf("Kode OTP Anda: %s. Jangan bagikan ke siapapun.", otp)

	// Prepare JSON payload
	payload := map[string]string{
		"token":   cfg.APIToken,
		"secret":  cfg.APISecret,
		"phone":   mobile,
		"message": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Create Request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp OTP: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("WhatsApp API error: %s", string(bodyBytes))
	}

	fmt.Printf("✅ WhatsApp OTP sent to %s. Response: %s\n", mobile, string(bodyBytes))
	return nil
}

// SendWhatsAppMessage sends a custom message via WhatsApp
func SendWhatsAppMessage(cfg WhatsAppConfig, mobile, message string) error {
	// Use URL from environment or fallback
	apiURL := os.Getenv("WA_URL_INVOICE")
	if apiURL == "" {
		apiURL = "https://apiwa.asitatech.co.id/api/sendMessageText"
	}

	// Prepare JSON payload
	payload := map[string]string{
		"token":   cfg.APIToken,
		"secret":  cfg.APISecret,
		"phone":   mobile,
		"message": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Create Request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("WhatsApp API error: %s", string(bodyBytes))
	}

	fmt.Printf("✅ WhatsApp message sent to %s\n", mobile)
	return nil
}

package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSettings fetches all system settings
func GetSettings(c *gin.Context) {
	var settings []models.Setting
	if err := config.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	// Group and format for easier UI consumption
	// Structure: { "payment": { "payment_qris_enabled": true, ... }, "general": { "app_name": "JalanJalanYuk" } }
	response := make(map[string]map[string]interface{})

	for _, s := range settings {
		if response[s.Group] == nil {
			response[s.Group] = make(map[string]interface{})
		}

		// Convert value based on type
		var val interface{}
		if s.Type == "bool" {
			b, _ := strconv.ParseBool(s.Value)
			val = b
		} else if s.Type == "number" {
			i, _ := strconv.Atoi(s.Value)
			val = i
		} else {
			val = s.Value
		}

		response[s.Group][s.Key] = val
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateSettings updates multiple settings by key
// Payload Example: { "payment_qris_enabled": false, "app_name": "NewName" }
func UpdateSettings(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	for key, val := range payload {
		var strVal string

		// Determine string representation
		switch v := val.(type) {
		case bool:
			strVal = strconv.FormatBool(v)
		case float64: // JSON numbers are float64 by default
			strVal = strconv.FormatFloat(v, 'f', -1, 64)
		case string:
			strVal = v
		default:
			continue // Skip unknown types
		}

		// Update in DB if key exists
		config.DB.Model(&models.Setting{}).Where("key = ?", key).Update("value", strVal)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

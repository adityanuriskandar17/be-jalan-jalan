package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadImage handles file uploads and returns the file URL
func UploadImage(c *gin.Context) {
	// 1. Get file from form-data
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 2. Validate file extension (optional)
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, JPEG, and PNG files are allowed"})
		return
	}

	// 3. Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	// Ensure uploads directory exists
	uploadPath := "./uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, 0755)
	}

	dst := filepath.Join(uploadPath, filename)

	// 4. Save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 5. Return URL
	// Construct full URL based on host
	protocol := "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", protocol, c.Request.Host)
	fileURL := fmt.Sprintf("%s/uploads/%s", baseURL, filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"data": gin.H{
			"url":      fileURL,
			"filename": filename,
		},
	})
}

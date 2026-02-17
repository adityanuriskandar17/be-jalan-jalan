package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- Statistics ---

// GetDestinationStats returns statistics for the destination dashboard
// GetDestinationStats returns statistics for the destination dashboard
func GetDestinationStats(c *gin.Context) {
	var total int64
	var active int64
	var inactive int64
	var avgRating float64

	// Total Destinations
	config.DB.Model(&models.Destination{}).Count(&total)

	// Active & Inactive
	config.DB.Model(&models.Destination{}).Where("is_active = ?", true).Count(&active)
	config.DB.Model(&models.Destination{}).Where("is_active = ?", false).Count(&inactive)

	// Average Rating
	config.DB.Model(&models.Destination{}).Select("COALESCE(AVG(rating_average), 0)").Scan(&avgRating)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_destinations":    total,
			"active_destinations":   active,
			"inactive_destinations": inactive,
			"average_rating":        avgRating,
		},
	})
}

// --- CRUD Destinations ---

// GetDestinations lists destinations with search and pagination
func GetDestinations(c *gin.Context) {
	var destinations []models.Destination
	// Start query with preloads
	query := config.DB.Model(&models.Destination{}).
		Preload("Categories").
		Preload("Images").
		Preload("Tickets") // Preload tickets

	// Search
	search := c.Query("search")
	if search != "" {
		searchStr := "%" + search + "%"
		// Filter by name, city, OR category name
		query = query.Joins("LEFT JOIN destination_categories ON destination_categories.destination_id = destinations.id").
			Joins("LEFT JOIN categories ON categories.id = destination_categories.category_id").
			Where("destinations.name ILIKE ? OR destinations.city ILIKE ? OR destinations.address ILIKE ? OR categories.name ILIKE ?", searchStr, searchStr, searchStr, searchStr).
			Group("destinations.id")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	if search != "" {
		config.DB.Model(&models.Destination{}).
			Joins("LEFT JOIN destination_categories ON destination_categories.destination_id = destinations.id").
			Joins("LEFT JOIN categories ON categories.id = destination_categories.category_id").
			Where("destinations.name ILIKE ? OR destinations.city ILIKE ? OR destinations.address ILIKE ? OR categories.name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").
			Distinct("destinations.id").
			Count(&total)
	} else {
		config.DB.Model(&models.Destination{}).Count(&total)
	}

	if err := query.Offset(offset).Limit(limit).Find(&destinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch destinations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destinations,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreateDestination adds a new destination
func CreateDestination(c *gin.Context) {
	var req models.Destination
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create destination"})
		return
	}

	// Reload data to populate association fields (Categories name, etc) for response
	config.DB.Preload("Categories").Preload("Tickets").Preload("Images").First(&req, req.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Destination created successfully", "data": req})
}

// UpdateDestination updates an existing destination
func UpdateDestination(c *gin.Context) {
	id := c.Param("id")
	var dest models.Destination
	if err := config.DB.First(&dest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}

	var req models.Destination
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if err := config.DB.Model(&dest).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update destination"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination updated", "data": dest})
}

// DeleteDestination removes a destination
func DeleteDestination(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Destination{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete destination"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination deleted"})
}

// ToggleDestinationStatus switches the active status of a destination
func ToggleDestinationStatus(c *gin.Context) {
	id := c.Param("id")
	var dest models.Destination
	if err := config.DB.First(&dest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}

	// Toggle status
	dest.IsActive = !dest.IsActive
	config.DB.Save(&dest)

	status := "inactive"
	if dest.IsActive {
		status = "active"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination is now " + status, "data": dest})
}

// --- CRUD Categories ---

// GetCategories lists all categories
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateCategory adds a new category
func CreateCategory(c *gin.Context) {
	var req models.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "data": req})
}

// UpdateCategory updates a category
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req models.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if err := config.DB.Model(&cat).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated", "data": cat})
}

// DeleteCategory deletes a category
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

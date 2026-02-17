package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- Statistics ---

// GetUserStats returns statistics for user management
func GetUserStats(c *gin.Context) {
	var total, active, suspended int64
	db := config.DB.Model(&models.User{})

	db.Count(&total)
	db.Where("is_active = ?", true).Count(&active)
	db.Where("is_active = ?", false).Count(&suspended)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_user": total,
			"active":     active,
			"suspended":  suspended,
		},
	})
}

// --- List Users ---

type UserListItem struct {
	ID           uint    `json:"id"`
	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number"`
	JoinDate     string  `json:"join_date"`
	TotalBooking int     `json:"total_booking"`
	TotalSpent   float64 `json:"total_spent"`
	Status       string  `json:"status"` // Active / Suspended
	IsActive     bool    `json:"is_active"`
}

// GetUsers lists all users with aggregation for bookings and spendings
func GetUsers(c *gin.Context) {
	var users []models.User
	query := config.DB.Model(&models.User{})

	// Search
	search := c.Query("search")
	if search != "" {
		query = query.Where("full_name ILIKE ? OR email ILIKE ? OR phone_number ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Transform to Response List (with aggregation)
	// Note: For high performance, use a raw SQL query or JOIN with Group By.
	// For now, loop and count is acceptable for admin panel scale.
	var userList []UserListItem
	for _, u := range users {
		var totalBooking int64
		var totalSpent float64

		// Count bookings (PAID transactions)
		config.DB.Model(&models.Order{}).Where("user_id = ? AND payment_status = ?", u.ID, "PAID").Count(&totalBooking)

		// Sum spent (PAID transactions)
		var sum *float64
		row := config.DB.Model(&models.Order{}).Where("user_id = ? AND payment_status = ?", u.ID, "PAID").Select("sum(total_amount)").Row()
		row.Scan(&sum)
		if sum != nil {
			totalSpent = *sum
		}

		status := "Suspended"
		if u.IsActive {
			status = "Active"
		}

		userList = append(userList, UserListItem{
			ID:           u.ID,
			FullName:     u.FullName,
			Email:        u.Email,
			PhoneNumber:  u.PhoneNumber,
			JoinDate:     u.CreatedAt.Format("02 Jan 2006"),
			TotalBooking: int(totalBooking),
			TotalSpent:   totalSpent,
			Status:       status,
			IsActive:     u.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userList,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// --- Status Management ---

// ToggleUserStatus suspends or activates a user
func ToggleUserStatus(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Toggle status
	user.IsActive = !user.IsActive
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	status := "Suspended"
	if user.IsActive {
		status = "Active"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User status updated",
		"data": gin.H{
			"id":        user.ID,
			"is_active": user.IsActive,
			"status":    status,
		},
	})
}

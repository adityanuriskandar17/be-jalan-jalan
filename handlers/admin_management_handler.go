package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminStats represents the summary cards
type AdminStats struct {
	TotalAdmin      int64 `json:"total_admin"`
	SuperAdmin      int64 `json:"super_admin"`
	Admin           int64 `json:"admin"`
	CustomerService int64 `json:"customer_service"`
}

// GetAdminStats returns counts of different admin roles
func GetAdminStats(c *gin.Context) {
	var stats AdminStats

	// Count all non-user roles
	config.DB.Model(&models.User{}).Where("role IN ?", []string{"admin", "super_admin", "customer_service"}).Count(&stats.TotalAdmin)

	config.DB.Model(&models.User{}).Where("role = ?", "super_admin").Count(&stats.SuperAdmin)
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.Admin)
	config.DB.Model(&models.User{}).Where("role = ?", "customer_service").Count(&stats.CustomerService)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetAdmins lists admin users with pagination and search
func GetAdmins(c *gin.Context) {
	var users []models.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	search := c.Query("search")
	role := c.Query("role")

	query := config.DB.Model(&models.User{}).Where("role IN ?", []string{"admin", "super_admin", "customer_service"})

	if search != "" {
		query = query.Where("full_name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	query.Count(&total)

	query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateAdmin creates a new admin/staff account
func CreateAdmin(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required,oneof=admin super_admin customer_service"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email exists
	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := models.User{
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully", "data": newUser})
}

// UpdateAdmin updates an existing admin account
func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"omitempty,min=8"`
		Role     string `json:"role" binding:"omitempty,oneof=admin super_admin customer_service"`
		IsActive *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Updates
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	} // Warning: Should check duplicate email if changed
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully", "data": user})
}

// DeleteAdmin deletes an admin account
func DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	// Prevent deleting self? (Optional, handled by frontend or policy)

	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete admin"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}

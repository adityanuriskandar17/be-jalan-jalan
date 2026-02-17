package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Structs for Requests ---

type PromoRequest struct {
	Code           string  `json:"code" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	DiscountType   string  `json:"discount_type" binding:"required"` // PERCENTAGE / FIXED
	DiscountValue  float64 `json:"discount_value" binding:"required"`
	MinTransaction float64 `json:"min_transaction"`
	MaxDiscount    float64 `json:"max_discount"`
	ValidFrom      string  `json:"valid_from" binding:"required"`  // YYYY-MM-DD
	ValidUntil     string  `json:"valid_until" binding:"required"` // YYYY-MM-DD
	Quota          int     `json:"quota"`                          // -1 unlimited
	IsActive       bool    `json:"is_active"`
}

// --- Statistics ---

// GetPromoStats returns promo statistics
func GetPromoStats(c *gin.Context) {
	var total, active int64

	config.DB.Model(&models.PromoCode{}).Count(&total)
	// Active = IsActive AND Not Expired
	config.DB.Model(&models.PromoCode{}).Where("is_active = ? AND valid_until >= ?", true, time.Now()).Count(&active)

	// Total Usage (Sum of UsedCount)
	var sum *int
	row := config.DB.Model(&models.PromoCode{}).Select("sum(used_count)").Row()
	row.Scan(&sum)

	totalUsage := 0
	if sum != nil {
		totalUsage = *sum
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_promo":  total,
			"active_promo": active,
			"total_usage":  totalUsage,
		},
	})
}

// --- CRUD ---

// GetPromos lists promos with search and pagination
func GetPromos(c *gin.Context) {
	var promos []models.PromoCode
	query := config.DB.Model(&models.PromoCode{})

	// Search
	search := c.Query("search")
	if search != "" {
		query = query.Where("code ILIKE ? OR title ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&promos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": promos,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetPromoDetail returns single promo detail
func GetPromoDetail(c *gin.Context) {
	id := c.Param("id")
	var promo models.PromoCode
	if err := config.DB.First(&promo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Promo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": promo})
}

// CreatePromo creates a new promo code
func CreatePromo(c *gin.Context) {
	var req PromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse Dates
	validFrom, err1 := time.Parse("2006-01-02", req.ValidFrom)
	validUntil, err2 := time.Parse("2006-01-02", req.ValidUntil)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	promo := models.PromoCode{
		Code:           req.Code,
		Title:          req.Title,
		Description:    req.Description,
		DiscountType:   req.DiscountType,
		DiscountValue:  req.DiscountValue,
		MinTransaction: req.MinTransaction,
		MaxDiscount:    req.MaxDiscount,
		ValidFrom:      validFrom,
		ValidUntil:     validUntil,
		Quota:          req.Quota,
		IsActive:       req.IsActive,
		UsedCount:      0,
	}

	if err := config.DB.Create(&promo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promo. Code might be duplicate."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Promo created", "data": promo})
}

// UpdatePromo updates existing promo
func UpdatePromo(c *gin.Context) {
	id := c.Param("id")
	var req PromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var promo models.PromoCode
	if err := config.DB.First(&promo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Promo not found"})
		return
	}

	// Parse Dates
	if req.ValidFrom != "" {
		validFrom, err := time.Parse("2006-01-02", req.ValidFrom)
		if err == nil {
			promo.ValidFrom = validFrom
		}
	}
	if req.ValidUntil != "" {
		validUntil, err := time.Parse("2006-01-02", req.ValidUntil)
		if err == nil {
			promo.ValidUntil = validUntil
		}
	}

	promo.Code = req.Code
	promo.Title = req.Title
	promo.Description = req.Description
	promo.DiscountType = req.DiscountType
	promo.DiscountValue = req.DiscountValue
	promo.MinTransaction = req.MinTransaction
	promo.MaxDiscount = req.MaxDiscount
	promo.Quota = req.Quota
	promo.IsActive = req.IsActive

	if err := config.DB.Save(&promo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update promo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promo updated", "data": promo})
}

// DeletePromo deletes a promo
func DeletePromo(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.PromoCode{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete promo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Promo deleted"})
}

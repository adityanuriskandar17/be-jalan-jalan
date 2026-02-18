package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- Statistics ---

// GetPaymentStats returns statistics for payment management
func GetPaymentStats(c *gin.Context) {
	var totalRevenue float64
	var pendingCount, totalTransactions int64
	var pendingAmount float64

	db := config.DB.Model(&models.Order{})

	// Total Transactions
	db.Count(&totalTransactions)

	// Pending Count & Amount
	db.Where("payment_status = ?", "UNPAID").Count(&pendingCount)

	var pendingSum *float64
	config.DB.Model(&models.Order{}).Where("payment_status = ?", "UNPAID").Select("sum(total_amount)").Row().Scan(&pendingSum)
	if pendingSum != nil {
		pendingAmount = *pendingSum
	}

	// Total Revenue (PAID transactions)
	var revenueSum *float64
	config.DB.Model(&models.Order{}).Where("payment_status = ?", "PAID").Select("sum(total_amount)").Row().Scan(&revenueSum)
	if revenueSum != nil {
		totalRevenue = *revenueSum
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_revenue":      totalRevenue,
			"pending_count":      pendingCount,
			"pending_amount":     pendingAmount,
			"total_transactions": totalTransactions,
		},
	})
}

// --- CRUD & Actions ---

// GetPayments lists transactions with filtering and pagination
func GetPayments(c *gin.Context) {
	var orders []models.Order
	query := config.DB.Model(&models.Order{}).Preload("User").Preload("Items.Ticket")

	// Search (ID, User Name, Payment Method)
	search := c.Query("search")
	if search != "" {
		// Use explicit JOIN to ensure table alias is correct.
		// Assuming table name is "users"
		query = query.Joins("JOIN users ON users.id = orders.user_id").Where(
			"orders.order_no ILIKE ? OR orders.visitor_name ILIKE ? OR orders.payment_method ILIKE ? OR users.full_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	// Filter Status (Payment Status)
	status := c.Query("status")
	if status != "" && status != "Semua Status" {
		query = query.Where("payment_status = ? OR status = ?", status, status) // Allow filtering by both statuses roughly
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetPaymentDetail returns single payment transaction detail
func GetPaymentDetail(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.Preload("User").Preload("Items.Ticket").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// ConfirmPayment approves a pending payment
func ConfirmPayment(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if order.PaymentStatus == "PAID" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is already PAID"})
		return
	}

	// Update Status
	updates := map[string]interface{}{
		"payment_status": "PAID",
		"status":         "CONFIRMED",
		"updated_at":     config.DB.NowFunc(),
	}

	if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment confirmed successfully", "data": order})
}

// RejectPayment rejects a pending payment
func RejectPayment(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if order.Status == "CANCELLED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is already CANCELLED"})
		return
	}

	// Update Status
	updates := map[string]interface{}{
		"payment_status": "CANCELLED", // Or FAILED if preferred
		"status":         "CANCELLED",
		"updated_at":     config.DB.NowFunc(),
	}

	if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment rejected/cancelled", "data": order})
}

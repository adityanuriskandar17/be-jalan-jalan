package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Statistics ---

// GetBookingStats returns statistics for booking management
func GetBookingStats(c *gin.Context) {
	var total, pending, confirmed, completed int64
	db := config.DB.Model(&models.Order{})

	db.Count(&total)
	db.Where("status = ?", "PENDING").Count(&pending)
	db.Where("status = ?", "CONFIRMED").Count(&confirmed)
	db.Where("status = ?", "COMPLETED").Count(&completed)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_booking": total,
			"pending":       pending,
			"confirmed":     confirmed,
			"completed":     completed,
		},
	})
}

// --- CRUD Bookings ---

// GetBookings lists all bookings with search and filter
func GetBookings(c *gin.Context) {
	var bookings []models.Order
	query := config.DB.Model(&models.Order{}).
		Preload("Items.Ticket.Destination"). // Deep preload
		Preload("User")                      // Preload relations

	// Search
	search := c.Query("search")
	if search != "" {
		// Search by OrderNo, VisitorName
		query = query.Where("order_no ILIKE ? OR visitor_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Filter Status
	status := c.Query("status")
	if status != "" && status != "Semua Status" {
		query = query.Where("status = ?", status)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetBookingDetail returns full details of a specific booking
func GetBookingDetail(c *gin.Context) {
	id := c.Param("id") // Can be ID or OrderNo (BK001)

	var order models.Order
	query := config.DB.Preload("Items.Ticket.Destination").Preload("User")

	// Check if input is numeric ID or string OrderNo
	// Try finding by ID first if numeric
	if _, err := strconv.Atoi(id); err == nil {
		if err := query.First(&order, id).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"data": order})
			return
		}
	}

	// Fallback to OrderNo search
	if err := query.Where("order_no = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// GetBookingETicket validates and returns data for E-Ticket generation
func GetBookingETicket(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := config.DB.Preload("Items.Ticket.Destination").Preload("User").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Validate Status
	// Allow case-insensitive check or strict check
	if order.PaymentStatus != "PAID" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-Ticket not available: Payment not settled (Unpaid)"})
		return
	}
	if order.Status != "CONFIRMED" && order.Status != "COMPLETED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-Ticket not available: Booking not confirmed"})
		return
	}

	// Prepare E-Ticket Data
	destinationName := "Unknown Destination"
	destinationCity := "Solo"
	if len(order.Items) > 0 && order.Items[0].Ticket.Destination != nil {
		destinationName = order.Items[0].Ticket.Destination.Name
		destinationCity = order.Items[0].Ticket.Destination.City
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "E-Ticket is valid",
		"data": gin.H{
			"order_no":         order.OrderNo,
			"qr_code_data":     order.OrderNo, // To be rendered as QR by Frontend
			"barcode_data":     order.OrderNo, // To be rendered as Barcode
			"visitor_name":     order.VisitorName,
			"visitor_email":    order.VisitorEmail,
			"visit_date":       order.BookingDate,
			"destination":      destinationName,
			"destination_city": destinationCity,
			"total_amount":     order.TotalAmount,
			"payment_status":   "LUNAS", // Display text
			"items":            order.Items,
		},
	})
}

// ManualBookingRequest struct for payload validation
type ManualBookingRequest struct {
	DestinationID uint   `json:"destination_id" binding:"required"`
	VisitorName   string `json:"visitor_name" binding:"required"`
	VisitorEmail  string `json:"visitor_email" binding:"required,email"`
	VisitorPhone  string `json:"visitor_phone" binding:"required"`
	VisitDate     string `json:"visit_date" binding:"required"` // YYYY-MM-DD
	Items         []struct {
		TicketID uint `json:"ticket_id" binding:"required"`
		Quantity int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,dive"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Status        string `json:"status" binding:"required"`         // CONFIRMED/PENDING
	PaymentStatus string `json:"payment_status" binding:"required"` // PAID/UNPAID
}

// CreateManualBooking handles manual booking creation by admin
func CreateManualBooking(c *gin.Context) {
	var req ManualBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse Date
	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (YYYY-MM-DD)"})
		return
	}

	// Find Destination (Validation)
	var dest models.Destination
	if err := config.DB.First(&dest, req.DestinationID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination not found"})
		return
	}

	// Calculate Total
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, itemReq := range req.Items {
		var ticket models.Ticket
		if err := config.DB.First(&ticket, itemReq.TicketID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ticket ID %d not found", itemReq.TicketID)})
			return
		}

		// Validation: Ticket belongs to Destination (or is Master Ticket linked?)
		if ticket.DestinationID != nil && *ticket.DestinationID != req.DestinationID {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ticket '%s' does not belong to selected destination", ticket.Name)})
			return
		}

		subTotal := ticket.Price * float64(itemReq.Quantity)
		totalAmount += subTotal

		orderItems = append(orderItems, models.OrderItem{
			TicketID:        ticket.ID,
			Quantity:        itemReq.Quantity,
			PriceAtPurchase: ticket.Price,
		})
	}

	// Create User Logic for Manual Booking
	var userID uint
	var user models.User
	if err := config.DB.Where("email = ?", req.VisitorEmail).First(&user).Error; err == nil {
		userID = user.ID
	} else {
		// Create new user implicitly
		newUser := models.User{
			FullName: req.VisitorName,
			Email:    req.VisitorEmail,
			Role:     "user",
			Password: "WalkInUser123!", // Dummy password
		}
		if err := config.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user record"})
			return
		}
		userID = newUser.ID
	}

	order := models.Order{
		UserID:        userID,
		OrderNo:       fmt.Sprintf("M-BK-%d", time.Now().Unix()), // Manual Booking Prefix
		VisitorName:   req.VisitorName,
		VisitorEmail:  req.VisitorEmail,
		VisitorPhone:  req.VisitorPhone,
		BookingDate:   visitDate,
		TotalAmount:   totalAmount,
		SubTotal:      totalAmount,
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		PaymentMethod: req.PaymentMethod,
		Items:         orderItems,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully", "data": order})
}

// UpdateBookingStatus updates status
func UpdateBookingStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status        string `json:"status"`
		PaymentStatus string `json:"payment_status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.PaymentStatus != "" {
		updates["payment_status"] = req.PaymentStatus
	}

	if err := config.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated", "data": order})
}

// DeleteBooking deletes a booking
func DeleteBooking(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Order{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted"})
}

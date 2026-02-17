package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- Statistics ---

// GetTicketStats returns statistics for ticket management
func GetTicketStats(c *gin.Context) {
	var total, reguler, paket, active int64

	config.DB.Model(&models.Ticket{}).Count(&total)
	config.DB.Model(&models.Ticket{}).Where("type = ? OR type = ?", "Reguler", "Tiket Reguler").Count(&reguler)
	config.DB.Model(&models.Ticket{}).Where("type LIKE ?", "%Paket%").Count(&paket)
	config.DB.Model(&models.Ticket{}).Where("is_active = ?", true).Count(&active)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_tickets":   total,
			"reguler_tickets": reguler,
			"paket_tickets":   paket,
			"active_tickets":  active,
		},
	})
}

// --- CRUD Tickets ---

// GetTickets lists all tickets with filtering
func GetTickets(c *gin.Context) {
	var tickets []models.Ticket
	query := config.DB.Model(&models.Ticket{})

	// Search
	search := c.Query("search")
	if search != "" {
		sql := "name ILIKE ? OR description ILIKE ?"
		arg := "%" + search + "%"
		query = query.Where(sql, arg, arg)
	}

	// Filter Type
	ticketType := c.Query("type") // Reguler, Paket
	if ticketType != "" {
		if ticketType == "Reguler" {
			query = query.Where("type = ? OR type = ?", "Reguler", "Tiket Reguler")
		} else if ticketType == "Paket" {
			query = query.Where("type LIKE ?", "%Paket%")
		} else {
			query = query.Where("type = ?", ticketType)
		}
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tickets,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreateTicket adds a new ticket (Master or linked to Destination)
func CreateTicket(c *gin.Context) {
	var req models.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ticket created successfully", "data": req})
}

// UpdateTicket updates an existing ticket
func UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	var ticket models.Ticket
	if err := config.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	var req models.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&ticket).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated", "data": ticket})
}

// DeleteTicket deletes a ticket
func DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Ticket{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
}

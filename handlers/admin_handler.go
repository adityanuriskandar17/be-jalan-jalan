package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	TotalRevenue        float64              `json:"total_revenue"`
	TotalBookings       int64                `json:"total_bookings"`
	ActiveUsers         int64                `json:"active_users"`
	PendingPayments     int64                `json:"pending_payments"`
	RecentBookings      []RecentBooking      `json:"recent_bookings"`
	PopularDestinations []PopularDestination `json:"popular_destinations"`
	CategoryStats       []CategoryStat       `json:"category_stats"`
	TicketStats         TicketStat           `json:"ticket_stats"`
}

type RecentBooking struct {
	ID          uint      `json:"id"`
	User        string    `json:"user"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
}

type PopularDestination struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	BookingCount int64  `json:"booking_count"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type TicketStat struct {
	ActiveCount int64 `json:"active_count"`
	TotalCount  int64 `json:"total_count"`
}

// GetDashboardStats retrieves statistics for the admin dashboard
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 1. Total Revenue (Only PAID orders)
	config.DB.Model(&models.Order{}).Where("status = ?", "PAID").Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalRevenue)

	// 2. Total Bookings
	config.DB.Model(&models.Order{}).Count(&stats.TotalBookings)

	// 3. Active Users (Role = 'user')
	config.DB.Model(&models.User{}).Where("role = ?", "user").Count(&stats.ActiveUsers)

	// 4. Pending Payments
	config.DB.Model(&models.Order{}).Where("status = ?", "PENDING").Count(&stats.PendingPayments)

	// 5. Recent Bookings (Limit 5)
	// We need to join with User and Items->Ticket to get ticket details
	var recentOrders []models.Order
	config.DB.Preload("User").Preload("Items.Ticket").Order("created_at desc").Limit(5).Find(&recentOrders)

	for _, order := range recentOrders {
		// Determine user name
		userName := order.VisitorName
		if userName == "" {
			userName = order.User.FullName
		}

		destName := "Unknown"
		if len(order.Items) > 0 {
			// Find the first destination from items
			var ticket models.Ticket
			if err := config.DB.First(&ticket, order.Items[0].TicketID).Error; err == nil {
				var dest models.Destination
				if err := config.DB.First(&dest, ticket.DestinationID).Error; err == nil {
					destName = dest.Name
				}
			}
		}

		stats.RecentBookings = append(stats.RecentBookings, RecentBooking{
			ID:          order.ID,
			User:        userName,
			Destination: destName,
			Date:        order.CreatedAt,
			Amount:      order.TotalAmount,
			Status:      order.Status,
		})
	}

	// 6. Popular Destinations
	// Count occurrence in OrderItems
	type PopDestResult struct {
		DestinationID uint
		Count         int64
	}
	var popResults []PopDestResult
	// Join OrderItem -> Ticket -> Destination
	// Group by DestinationID
	// This query depends on exact schema. Simplified approach:
	config.DB.Table("order_items").
		Joins("JOIN tickets ON tickets.id = order_items.ticket_id").
		Select("tickets.destination_id, count(*) as count").
		Group("tickets.destination_id").
		Order("count desc").
		Limit(5).
		Scan(&popResults)

	for _, res := range popResults {
		var dest models.Destination
		config.DB.First(&dest, res.DestinationID)
		stats.PopularDestinations = append(stats.PopularDestinations, PopularDestination{
			ID:           dest.ID,
			Name:         dest.Name,
			BookingCount: res.Count,
		})
	}

	// 7. Category Stats
	// Count destinations per category
	// Need to query destination_categories pivot table
	type CatResult struct {
		CategoryName string
		Count        int64
	}
	// This is a complex join. Let's simplify: Count destinations per category
	// SELECT c.name, count(d.id) FROM categories c JOIN destination_categories dc ON c.id = dc.category_id JOIN destinations d ON d.id = dc.destination_id GROUP BY c.name
	rows, err := config.DB.Raw(`
		SELECT c.name, count(dc.destination_id) 
		FROM categories c 
		JOIN destination_categories dc ON c.id = dc.category_id 
		GROUP BY c.name
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var count int64
			rows.Scan(&name, &count)
			stats.CategoryStats = append(stats.CategoryStats, CategoryStat{Category: name, Count: count})
		}
	}

	// 8. Ticket Stats
	config.DB.Model(&models.Ticket{}).Count(&stats.TicketStats.TotalCount)
	// Assuming no 'active' status on tickets yet, so Active = Total for now
	stats.TicketStats.ActiveCount = stats.TicketStats.TotalCount

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard stats retrieved successfully",
		"data":    stats,
	})
}

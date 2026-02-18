package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// --- Structs ---
type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type TopDestination struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	BookingCount int64   `json:"booking_count"`
	Revenue      float64 `json:"revenue"`
}

type PaymentMethodStats struct {
	Method string  `json:"method"`
	Count  int64   `json:"count"`
	Total  float64 `json:"total"`
}

type CategoryPerformance struct {
	Category string `json:"category"`
	Count    int64  `json:"count"` // Number of destinations
}

type KeyMetrics struct {
	AvgOrderValue  float64 `json:"avg_order_value"`
	ConversionRate float64 `json:"conversion_rate"`
	RepeatCustomer float64 `json:"repeat_customer"` // Percentage
	AvgLeadTime    int     `json:"avg_lead_time"`   // Days
	TotalRevenue   float64 `json:"total_revenue"`
}

// --- Handler ---

// GetDashboardReports returns aggregated data for analytics
func GetDashboardReports(c *gin.Context) {
	period := c.DefaultQuery("period", "6_months") // 7_days, 30_days, 3_months, 6_months, 1_year
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error
	now := time.Now()

	// 1. Determine Date Range
	if startDateStr != "" && endDateStr != "" {
		// Custom Range
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format (YYYY-MM-DD)"})
			return
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (YYYY-MM-DD)"})
			return
		}
		// Set endDate to end of that day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		period = "custom"
	} else {
		// Presets
		endDate = now
		switch period {
		case "7_days":
			startDate = now.AddDate(0, 0, -7)
		case "30_days":
			startDate = now.AddDate(0, 0, -30)
		case "3_months":
			startDate = now.AddDate(0, -3, 0)
		case "6_months":
			startDate = now.AddDate(0, -6, 0)
		case "1_year":
			startDate = now.AddDate(-1, 0, 0)
		default:
			startDate = now.AddDate(0, -6, 0) // Default 6 months
		}
	}

	// 2. Fetch Data (Filtered by Range)
	var orders []models.Order
	config.DB.Where("payment_status = ? AND created_at >= ? AND created_at <= ?", "PAID", startDate, endDate).Find(&orders)

	// 3. Prepare Chart Data (Dynamic Aggregation)
	var chartData []MonthlyRevenue
	durationDays := endDate.Sub(startDate).Hours() / 24

	if durationDays <= 32 { // Show Daily if range is ~1 month or less
		// Daily Aggregation
		dailyMap := make(map[string]float64)
		for _, o := range orders {
			day := o.CreatedAt.Format("02 Jan")
			dailyMap[day] += o.TotalAmount
		}

		// Fill missing days
		daysToFill := int(durationDays)
		if daysToFill < 0 {
			daysToFill = 0
		} // Safety

		// Reconstruct timeline from Start to End
		// To ensure correct order, iterate from startDate to endDate
		current := startDate
		for !current.After(endDate) {
			dStr := current.Format("02 Jan")
			rev := 0.0
			// Look up in map (optimization: map is faster than iterating orders again)
			if val, ok := dailyMap[dStr]; ok {
				rev = val
			}
			chartData = append(chartData, MonthlyRevenue{Month: dStr, Revenue: rev})
			current = current.AddDate(0, 0, 1) // Next day
		}

	} else {
		// Monthly Aggregation
		revenueMap := make(map[string]float64)
		for _, o := range orders {
			mStr := o.CreatedAt.Format("Jan 2006")
			revenueMap[mStr] += o.TotalAmount
		}

		// Reconstruct timeline from Start to End (Month by Month)
		current := startDate
		// Normalize to first of month to avoid skipping if started on 31st
		current = time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())

		for !current.After(endDate) {
			mStr := current.Format("Jan 2006")
			rev := 0.0
			if val, ok := revenueMap[mStr]; ok {
				rev = val
			}
			chartData = append(chartData, MonthlyRevenue{Month: mStr, Revenue: rev})

			// Add 1 month
			current = current.AddDate(0, 1, 0)
		}
	}

	// 4. Top Destinations (Filtered)
	var topDestinations []TopDestination
	rows, _ := config.DB.Table("order_items").
		Select("destinations.id, destinations.name, count(order_items.id) as booking_count, sum(order_items.price_at_purchase * order_items.quantity) as revenue").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN tickets ON tickets.id = order_items.ticket_id").
		Joins("JOIN destinations ON destinations.id = tickets.destination_id").
		Where("orders.payment_status = ? AND orders.created_at >= ? AND orders.created_at <= ?", "PAID", startDate, endDate).
		Group("destinations.id, destinations.name").
		Order("revenue desc").
		Limit(5).
		Rows()

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var td TopDestination
			rows.Scan(&td.ID, &td.Name, &td.BookingCount, &td.Revenue)
			topDestinations = append(topDestinations, td)
		}
	}

	// 5. Payment Methods (Filtered)
	var paymentMethods []PaymentMethodStats
	config.DB.Model(&models.Order{}).
		Select("payment_method as method, count(*) as count, sum(total_amount) as total").
		Where("payment_status = ? AND created_at >= ? AND created_at <= ?", "PAID", startDate, endDate).
		Group("payment_method").
		Scan(&paymentMethods)

	// 6. Metrics (Filtered)
	var totalOrders int64
	config.DB.Model(&models.Order{}).Where("created_at >= ? AND created_at <= ?", startDate, endDate).Count(&totalOrders)

	var paidOrders int64 = int64(len(orders))
	var totalRevenue float64
	for _, o := range orders {
		totalRevenue += o.TotalAmount
	}

	avgOrderValue := 0.0
	if paidOrders > 0 {
		avgOrderValue = totalRevenue / float64(paidOrders)
	}

	conversion := 0.0
	if totalOrders > 0 {
		conversion = (float64(paidOrders) / float64(totalOrders)) * 100
	}

	// Repeat Customer (In this period)
	var repeatCount int64
	config.DB.Model(&models.Order{}).
		Select("user_id").
		Where("payment_status = ? AND created_at >= ?", "PAID", startDate).
		Group("user_id").
		Having("count(id) > 1").
		Count(&repeatCount)

	var totalPaidUsers int64
	config.DB.Model(&models.Order{}).Where("payment_status = ? AND created_at >= ?", "PAID", startDate).Distinct("user_id").Count(&totalPaidUsers)

	repeatRate := 0.0
	if totalPaidUsers > 0 {
		repeatRate = (float64(repeatCount) / float64(totalPaidUsers)) * 100
	}

	// Avg Lead Time
	var leadTimeTotal float64
	var leadTimeCount int
	for _, o := range orders {
		diff := o.BookingDate.Sub(o.CreatedAt).Hours() / 24
		if diff > 0 {
			leadTimeTotal += diff
			leadTimeCount++
		}
	}
	avgLeadTime := 0
	if leadTimeCount > 0 {
		avgLeadTime = int(leadTimeTotal / float64(leadTimeCount))
	}

	metrics := KeyMetrics{
		AvgOrderValue:  avgOrderValue,
		ConversionRate: conversion,
		RepeatCustomer: repeatRate,
		AvgLeadTime:    avgLeadTime,
		TotalRevenue:   totalRevenue,
	}

	// 5. Category Performance (Filtered by Date via join?)
	// To strictly filter by date, we only count destinations that had BOOKINGS in this period?
	// Or just listing count? Usually "Performance" implies bookings/sales.
	// But the struct only has "Count". If it means "Number of Destinations available", date doesn't matter much.
	// If it means "Number of Bookings per Category", we need to change logic.
	// Based on "Category Performance", sales/bookings is more relevant.
	// Let's update struct to be more useful? Or stick to previous logic?
	// Previous: "count(destination_categories.destination_id)" -> This is just how many dests exist in category. Static.
	// If user wants "Performance", it should probably be Sales.
	// However, to avoid breaking frontend expectation of "Count" meaning "Inventory Count", I will keep it static logic for now.
	// If specific request comes for "Sales per Category", we change it.
	// The prompt implies "Analisis performa", so Sales makes more sense, BUT variable name "count" usually implies quantity of items.
	// Let's Keep the current "Inventory Count" logic as it might be 'Number of Destinations in this Category'.

	var catPerf []CategoryPerformance
	cRows, _ := config.DB.Table("categories").
		Select("categories.name, count(destination_categories.destination_id) as count").
		Joins("LEFT JOIN destination_categories ON destination_categories.category_id = categories.id").
		Group("categories.id, categories.name").
		Rows()

	if cRows != nil {
		defer cRows.Close()
		for cRows.Next() {
			var cp CategoryPerformance
			cRows.Scan(&cp.Category, &cp.Count)
			catPerf = append(catPerf, cp)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"period":           period,
		"revenue_chart":    chartData,
		"top_destinations": topDestinations,
		"payment_methods":  paymentMethods,
		"metrics":          metrics,
		"category_perf":    catPerf,
	})
}

// ExportReports generates downloadable report
func ExportReports(c *gin.Context) {
	format := c.Query("format")

	// Fetch Data (Duplicate logic or reuse? For simplicity, fetch all paid orders)
	var orders []models.Order
	config.DB.Preload("User").Where("payment_status = ?", "PAID").Find(&orders)

	switch format {
	case "csv":
		c.Header("Content-Disposition", "attachment; filename=report.csv")
		c.Header("Content-Type", "text/csv")
		writer := csv.NewWriter(c.Writer)
		writer.Write([]string{"Order No", "Date", "Customer", "Amount", "Payment Method", "Status"})
		for _, o := range orders {
			writer.Write([]string{
				o.OrderNo,
				o.CreatedAt.Format("2006-01-02 15:04"),
				o.VisitorName,
				fmt.Sprintf("%.2f", o.TotalAmount),
				o.PaymentMethod,
				o.Status,
			})
		}
		writer.Flush()

	case "xlsx":
		f := excelize.NewFile()
		sheet := "Sheet1"
		headers := []string{"Order No", "Date", "Customer", "Amount", "Payment Method", "Status"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
		}
		for i, o := range orders {
			row := i + 2
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), o.OrderNo)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), o.CreatedAt.Format("2006-01-02 15:04"))
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), o.VisitorName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), o.TotalAmount)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), o.PaymentMethod)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), o.Status)
		}

		c.Header("Content-Disposition", "attachment; filename=report.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		f.Write(c.Writer)

	case "json":
		c.Header("Content-Disposition", "attachment; filename=report.json")
		c.JSON(http.StatusOK, orders)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
	}
}

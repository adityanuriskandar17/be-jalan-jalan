package database

import (
	"be-jalan/models"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedTravelData(db *gorm.DB) {
	// 1. Seed Categories
	categories := []models.Category{
		{Name: "Sejarah & Budaya", IconURL: "https://example.com/icon-culture.png"},
		{Name: "Alam", IconURL: "https://example.com/icon-nature.png"},
		{Name: "Kuliner", IconURL: "https://example.com/icon-food.png"},
		{Name: "Keluarga", IconURL: "https://example.com/icon-family.png"},
		{Name: "Belanja", IconURL: "https://example.com/icon-shop.png"},
	}
	for i, cat := range categories {
		db.FirstOrCreate(&categories[i], models.Category{Name: cat.Name})
	}

	// 2. Seed Destinations
	destinations := []models.Destination{
		{
			Name:           "Keraton Kasunanan Surakarta",
			Description:    "Istana resmi Kesultanan Surakarta Hadiningrat.",
			Address:        "Jl. Kamandungan, Baluwarti, Pasar Kliwon, Surakarta",
			City:           "Surakarta",
			PriceStartFrom: 15000,
			RatingAverage:  4.8,
			IsPopular:      true,
			Categories:     []models.Category{categories[0]}, // Sejarah
		},
		{
			Name:           "Taman Balekambang",
			Description:    "Taman peninggalan Mangkunegara VII yang asri.",
			Address:        "Jl. Balekambang, Manahan, Banjarsari, Surakarta",
			City:           "Surakarta",
			PriceStartFrom: 0,
			RatingAverage:  4.7,
			IsPopular:      true,
			Categories:     []models.Category{categories[1], categories[3]}, // Alam, Keluarga
		},
		{
			Name:           "Solo Safari",
			Description:    "Kebun binatang modern di Solo.",
			Address:        "Jl. Ir. Sutami No.109, Jebres, Surakarta",
			City:           "Surakarta",
			PriceStartFrom: 45000,
			RatingAverage:  4.5,
			IsPopular:      true,
			Categories:     []models.Category{categories[1], categories[3]}, // Alam, Keluarga
		},
	}

	for i, dest := range destinations {
		if err := db.Where("name = ?", dest.Name).FirstOrCreate(&destinations[i]).Error; err != nil {
			log.Printf("Failed to seed destination %s: %v", dest.Name, err)
		}
	}

	// 3. Seed Tickets
	tickets := []models.Ticket{
		{
			DestinationID: &destinations[0].ID, // Keraton
			Name:          "Tiket Masuk Dewasa",
			Type:          "Reguler",
			Price:         15000,
		},
		{
			DestinationID: &destinations[0].ID, // Keraton
			Name:          "Tiket Masuk Anak",
			Type:          "Reguler",
			Price:         10000,
		},
		{
			DestinationID: &destinations[1].ID, // Balekambang
			Name:          "Tiket Masuk",
			Type:          "Reguler",
			Price:         0,
		},
		{
			DestinationID: &destinations[2].ID, // Safari
			Name:          "Tiket Reguler",
			Type:          "Reguler",
			Price:         45000,
		},
		{
			DestinationID: &destinations[0].ID, // Keraton
			Name:          "Paket Foto Adat",
			Type:          "Paket",
			Price:         50000,
			OriginalPrice: 75000,
		},
		{
			DestinationID: &destinations[2].ID, // Safari
			Name:          "Paket Safari Feeding",
			Type:          "Paket",
			Price:         75000,
			OriginalPrice: 100000,
		},
	}
	for i, ticket := range tickets {
		var existingTicket models.Ticket
		// Check if ticket exists
		if err := db.Where("destination_id = ? AND name = ?", ticket.DestinationID, ticket.Name).First(&existingTicket).Error; err == nil {
			// Exists: Update it (especially Type)
			db.Model(&existingTicket).Updates(map[string]interface{}{
				"type":           ticket.Type,
				"price":          ticket.Price,
				"original_price": ticket.OriginalPrice,
			})
		} else {
			// Not exists: Create it
			db.Create(&tickets[i])
		}
	}

	// 4. Seed Dummy Users (for bookings)
	dummyUser := models.User{
		FullName: "Budi Santoso",
		Email:    "budi@example.com",
		Role:     "user",
	}
	db.Where("email = ?", dummyUser.Email).FirstOrCreate(&dummyUser)

	// 5. Seed Dummy Orders (Transactions)
	// Order 1: Success
	order1 := models.Order{
		UserID:        dummyUser.ID,
		OrderNo:       fmt.Sprintf("ORD-%d", time.Now().Unix()),
		TicketCode:    fmt.Sprintf("TKT-%d%s", time.Now().UnixNano()/1000, "A"),
		VisitorName:   "Budi Santoso",
		Status:        "CONFIRMED",
		PaymentStatus: "PAID",
		TotalAmount:   30000, // 2 Tiket Keraton
		BookingDate:   time.Now(),
		PaymentMethod: "QRIS",
		Items: []models.OrderItem{
			{TicketID: tickets[0].ID, Quantity: 2, PriceAtPurchase: 15000},
		},
	}

	// Check if order exists (simple check to avoid duplicate dummy data on every restart)
	var existingOrder models.Order
	if err := db.Where("visitor_name = ?", "Budi Santoso").First(&existingOrder).Error; err == nil {
		// FORCE UPDATE STATUS to ensure we have a PAID booking for testing
		db.Model(&existingOrder).Updates(map[string]interface{}{
			"status":         "CONFIRMED",
			"payment_status": "PAID",
			"ticket_code":    fmt.Sprintf("TKT-%d%s", time.Now().UnixNano()/1000, "U"),
		})
	} else {
		db.Create(&order1)

		// Order 2: Pending
		order2 := models.Order{
			UserID:        dummyUser.ID,
			OrderNo:       fmt.Sprintf("ORD-%d", time.Now().Unix()+1),
			TicketCode:    fmt.Sprintf("TKT-%d%s", time.Now().UnixNano()/1000, "B"),
			VisitorName:   "Siti Rahayu",
			Status:        "PENDING",
			PaymentStatus: "UNPAID",
			TotalAmount:   135000,
			BookingDate:   time.Now().AddDate(0, 0, 1),
			Items: []models.OrderItem{
				{TicketID: tickets[3].ID, Quantity: 3, PriceAtPurchase: 45000}, // 3 Tiket Safari
			},
		}
		db.Create(&order2)

		log.Println("✅ Seeded Travel Data & Dummy Orders")
	}
}

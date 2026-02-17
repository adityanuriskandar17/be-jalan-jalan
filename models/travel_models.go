package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a travel category (e.g. Museum, Nature, Culinary)
type Category struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"type:varchar(100);not null" json:"name"`
	IconURL      string        `gorm:"type:varchar(255)" json:"icon_url"`
	Destinations []Destination `gorm:"many2many:destination_categories;" json:"destinations,omitempty"`
}

// Destination represents a tourist spot
type Destination struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	Name           string  `gorm:"type:varchar(255);not null" json:"name"`
	Description    string  `gorm:"type:text" json:"description"`
	Address        string  `gorm:"type:text" json:"address"`
	City           string  `gorm:"type:varchar(100);default:'Solo'" json:"city"`
	Facilities     string  `gorm:"type:text" json:"facilities"` // Strings data (JSON/Comma separated)
	Latitude       float64 `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude      float64 `gorm:"type:decimal(11,8)" json:"longitude"`
	PriceStartFrom float64 `gorm:"type:decimal(10,2)" json:"price_start_from"`
	RatingAverage  float64 `gorm:"type:decimal(3,2);default:0" json:"rating_average"`
	ReviewCount    int     `gorm:"default:0" json:"review_count"`
	OpenTime       string  `gorm:"type:varchar(10)" json:"open_time"`  // Format "08:00"
	CloseTime      string  `gorm:"type:varchar(10)" json:"close_time"` // Format "17:00"

	// Relations
	Images     []DestinationImage `gorm:"foreignKey:DestinationID" json:"images"`
	Categories []Category         `gorm:"many2many:destination_categories;" json:"categories"`
	Tickets    []Ticket           `gorm:"foreignKey:DestinationID" json:"tickets"`
	Reviews    []Review           `gorm:"foreignKey:DestinationID" json:"reviews"`

	IsPopular bool `gorm:"default:false" json:"is_popular"` // For "Popular" badge/tab
	IsActive  bool `gorm:"default:true" json:"is_active"`   // For Active/Inactive toggle

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// DestinationImage represents a photo gallery for a destination
type DestinationImage struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	DestinationID uint   `gorm:"index;not null" json:"destination_id"`
	ImageURL      string `gorm:"type:varchar(255);not null" json:"image_url"`
	IsMain        bool   `gorm:"default:false" json:"is_main"` // Main thumbnail image
}

// Ticket represents purchaseable entry tickets
type Ticket struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	DestinationID *uint        `gorm:"index" json:"destination_id"`                           // Nullable (nil = Master Ticket)
	Destination   *Destination `gorm:"foreignKey:DestinationID" json:"destination,omitempty"` // Relation
	Name          string       `gorm:"type:varchar(255);not null" json:"name"`                // e.g., "Adult Ticket"
	Type          string       `gorm:"type:varchar(100)" json:"type"`                         // e.g. "Reguler", "Paket"
	Description   string       `gorm:"type:text" json:"description"`
	Price         float64      `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64      `gorm:"type:decimal(10,2)" json:"original_price"` // New: Strikethrough price
	Stock         int          `gorm:"default:-1" json:"stock"`                  // -1 means unlimited
	IsActive      bool         `gorm:"default:true" json:"is_active"`
	ValidityDate  *time.Time   `json:"validity_date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Review represents user ratings and comments
type Review struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	UserID        uint   `gorm:"index;not null" json:"user_id"`
	DestinationID uint   `gorm:"index;not null" json:"destination_id"`
	Rating        int    `gorm:"check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment       string `gorm:"type:text" json:"comment"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// Favorite represents user wishlist
type Favorite struct {
	UserID        uint `gorm:"primaryKey"`
	DestinationID uint `gorm:"primaryKey"`
	CreatedAt     time.Time

	Destination Destination `gorm:"foreignKey:DestinationID" json:"destination,omitempty"`
}

// Notification represents user alerts
type Notification struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Message string `gorm:"type:text;not null" json:"message"`
	Type    string `gorm:"type:varchar(50);not null" json:"type"` // INFO, PROMO, TRANSACTION
	IsRead  bool   `gorm:"default:false" json:"is_read"`

	CreatedAt time.Time `json:"created_at"`
}

// PromoCode represents discount codes
type PromoCode struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Code           string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Title          string    `gorm:"type:varchar(255);not null" json:"title"` // e.g. "Hemat Hingga 50%"
	Description    string    `gorm:"type:text" json:"description"`
	ImageURL       string    `gorm:"type:varchar(255)" json:"image_url"` // Banner image
	TermsCondition string    `gorm:"type:text" json:"terms_condition"`
	DiscountType   string    `gorm:"type:varchar(20);not null" json:"discount_type"` // PERCENTAGE or FIXED
	DiscountValue  float64   `gorm:"type:decimal(10,2);not null" json:"discount_value"`
	MaxDiscount    float64   `gorm:"type:decimal(10,2)" json:"max_discount"`
	MinTransaction float64   `gorm:"type:decimal(10,2);default:0" json:"min_transaction"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidUntil     time.Time `json:"valid_until"`
	Quota          int       `json:"quota"` // -1 for unlimited
	UsedCount      int       `json:"used_count" gorm:"default:0"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Order represents a transaction
type Order struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	PromoCodeID *uint  `gorm:"index" json:"promo_code_id"` // Nullable
	OrderNo     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_no"`
	TicketCode  string `gorm:"type:varchar(50);uniqueIndex" json:"ticket_code"` // Secure Random Code for QR

	// Snapshot Visitor Info (Data Pemesan)
	VisitorName  string `gorm:"type:varchar(100)" json:"visitor_name"`
	VisitorEmail string `gorm:"type:varchar(100)" json:"visitor_email"`
	VisitorPhone string `gorm:"type:varchar(20)" json:"visitor_phone"`

	// Payment Breakdown
	SubTotal       float64 `gorm:"type:decimal(10,2);not null" json:"sub_total"`
	ServiceFee     float64 `gorm:"type:decimal(10,2);default:0" json:"service_fee"`
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	TotalAmount    float64 `gorm:"type:decimal(10,2);not null" json:"total_amount"`

	// Payment Status & Details
	Status           string `gorm:"type:varchar(20);default:'PENDING'" json:"status"`        // Booking Status: PENDING, CONFIRMED, COMPLETED, CANCELLED
	PaymentStatus    string `gorm:"type:varchar(20);default:'UNPAID'" json:"payment_status"` // Pay Status: UNPAID, PAID
	PaymentMethod    string `gorm:"type:varchar(50)" json:"payment_method"`                  // QRIS, VA, E-WALLET
	PaymentChannel   string `gorm:"type:varchar(50)" json:"payment_channel"`                 // GoPay, BCA, Mandiri
	PaymentReference string `gorm:"type:varchar(255)" json:"payment_reference"`              // VA Number / QR String
	PaymentURL       string `gorm:"type:text" json:"payment_url"`                            // Deep link / Invoice URL

	PaymentDeadline time.Time  `json:"payment_deadline"` // Expiry time for payment
	PaymentDate     *time.Time `json:"payment_date"`
	BookingDate     time.Time  `json:"booking_date"` // Visit date

	// Relations
	Items     []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	User      User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PromoCode PromoCode   `gorm:"foreignKey:PromoCodeID" json:"promo_code,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderItem represents line items in an order
type OrderItem struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	OrderID         uint    `gorm:"index;not null" json:"order_id"`
	TicketID        uint    `gorm:"index;not null" json:"ticket_id"`
	Quantity        int     `gorm:"not null" json:"quantity"`
	PriceAtPurchase float64 `gorm:"type:decimal(10,2);not null" json:"price_at_purchase"`

	Ticket Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatRoom represents a conversation session
type ChatRoom struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`         // Customer
	AdminID       *uint     `gorm:"index" json:"admin_id"`        // Assigned Admin/Operator
	Status        string    `gorm:"default:'OPEN'" json:"status"` // OPEN, CLOSED, ASSIGNED
	LastMessageAt time.Time `json:"last_message_at"`

	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Admin    *User     `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Messages []Message `gorm:"foreignKey:ChatRoomID" json:"messages,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Message represents individual chat messages
type Message struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ChatRoomID uint   `gorm:"index" json:"chat_room_id"`
	SenderID   uint   `json:"sender_id"`   // User ID or Admin ID
	SenderType string `json:"sender_type"` // USER, ADMIN, SYSTEM
	Content    string `gorm:"type:text" json:"content"`
	IsRead     bool   `gorm:"default:false" json:"is_read"`

	CreatedAt time.Time `json:"created_at"`
}

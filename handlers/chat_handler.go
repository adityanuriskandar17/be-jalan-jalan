package handlers

import (
	"be-jalan/config"
	"be-jalan/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// --- WebSocket Hub (In-Memory) ---

type Client struct {
	ID   uint   // User or Admin ID
	Role string // "user" or "admin"
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	clients    map[string]*Client // Key: "role_id" (e.g. "user_1", "admin_2")
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.Mutex
}

var ChatHub = Hub{
	clients:    make(map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan []byte),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			key := client.Role + "_" + strconv.Itoa(int(client.ID))
			h.clients[key] = client
			h.mutex.Unlock()
			log.Printf("Client registered: %s", key)

		case client := <-h.unregister:
			h.mutex.Lock()
			key := client.Role + "_" + strconv.Itoa(int(client.ID))
			if _, ok := h.clients[key]; ok {
				delete(h.clients, key)
				close(client.Send)
			}
			h.mutex.Unlock()
			log.Printf("Client unregistered: %s", key)

		case <-h.broadcast:
			// Broadcast logic placeholder
			// currently handled in ReadPump directly
		}
	}
}

// Websocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for dev
	},
}

// WS Message Payload
type WSMessage struct {
	Type       string      `json:"type"` // "message", "typing", "read"
	ChatRoomID uint        `json:"chat_room_id"`
	SenderID   uint        `json:"sender_id"`
	SenderRole string      `json:"sender_role"` // "user" or "admin"
	Content    string      `json:"content"`
	CreatedAt  string      `json:"created_at"`
	Payload    interface{} `json:"payload,omitempty"`
}

// ServeWS handles WebSocket requests
func ServeWS(c *gin.Context) {
	// 1. Upgrade HTTP to WS
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WS Upgrade Error:", err)
		return
	}

	// 2. Identify User
	// Assume AuthMiddleware has set User/Admin ID or pass via Query
	// For simplicity, let's take query params: ?id=1&role=admin
	idStr := c.Query("id")
	role := c.Query("role") // "user" or "admin"

	id, _ := strconv.Atoi(idStr)
	if id == 0 || role == "" {
		conn.Close()
		return
	}

	client := &Client{ID: uint(id), Role: role, Conn: conn, Send: make(chan []byte, 256)}
	ChatHub.register <- client

	// 3. Goroutines for Pump
	go client.WritePump()
	go client.ReadPump()
}

func (c *Client) ReadPump() {
	defer func() {
		ChatHub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse Message
		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// Process Message
		if msg.Type == "message" {
			// Save to DB
			dbMsg := models.Message{
				ChatRoomID: msg.ChatRoomID,
				SenderID:   c.ID,
				SenderType: c.Role, // user or admin
				Content:    msg.Content,
				CreatedAt:  time.Now(),
			}
			config.DB.Create(&dbMsg) // Async save

			// Update ChatRoom LastMessage
			config.DB.Model(&models.ChatRoom{}).Where("id = ?", msg.ChatRoomID).
				Update("last_message_at", time.Now())

			// Route Message
			// Find Recipient in specific ChatRoom logic
			// 1. Find ChatRoom to get UserID and AdminID
			var room models.ChatRoom
			config.DB.First(&room, msg.ChatRoomID)

			var recipientKey string
			if c.Role == "user" {
				// Send to Assigned Admin (if any) or Broadcast to All Online Admins?
				// Usually to Assigned Admin
				if room.AdminID != nil {
					recipientKey = "admin_" + strconv.Itoa(int(*room.AdminID))
				} else {
					// No admin assigned yet, maybe broadcast to all admins?
					// For simple demo, let's just log it or try to send to any online admin
					// Or frontend polls the list.
				}
			} else {
				// Admin sending to User
				recipientKey = "user_" + strconv.Itoa(int(room.UserID))
			}

			// Send back to Sender (Confirmation/Echo)
			// senderKey := c.Role + "_" + strconv.Itoa(int(c.ID))
			// if client, ok := ChatHub.clients[senderKey]; ok {
			// 	client.Send <- message
			// }

			// Send to Recipient if online
			ChatHub.mutex.Lock()
			if recipient, ok := ChatHub.clients[recipientKey]; ok {
				recipient.Send <- message
			}
			ChatHub.mutex.Unlock()
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// --- REST API Handlers ---

// GetChatRooms (Admin sidebar)
func GetChatRooms(c *gin.Context) {
	var rooms []models.ChatRoom
	status := c.Query("status") // OPEN, RESOLVED

	query := config.DB.Preload("User").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(1) // Preview last message
	})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("last_message_at desc").Find(&rooms)
	c.JSON(http.StatusOK, gin.H{"data": rooms})
}

// GetChatMessages (Chat window)
func GetChatMessages(c *gin.Context) {
	roomID := c.Param("id")
	var messages []models.Message

	if err := config.DB.Where("chat_room_id = ?", roomID).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// InitChat (User starts chat)
func InitChat(c *gin.Context) {
	// Assume User Auth
	// For demo, receive user_id in body
	var req struct {
		UserID uint `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check existing OPEN chat
	var room models.ChatRoom
	if err := config.DB.Where("user_id = ? AND status = ?", req.UserID, "OPEN").First(&room).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Continuing existing chat", "data": room})
		return
	}

	// Create New
	newRoom := models.ChatRoom{
		UserID:        req.UserID,
		Status:        "OPEN",
		LastMessageAt: time.Now(),
	}
	config.DB.Create(&newRoom)

	// Create Welcome Message
	welcome := models.Message{
		ChatRoomID: newRoom.ID,
		SenderID:   0, // System
		SenderType: "SYSTEM",
		Content:    "Selamat datang di Live Chat Support. Ada yang bisa kami bantu?",
		CreatedAt:  time.Now(),
	}
	config.DB.Create(&welcome)

	c.JSON(http.StatusCreated, gin.H{"message": "Chat initialized", "data": newRoom})
}

// JoinChat (Admin joins/assigns self)
func JoinChat(c *gin.Context) {
	roomID := c.Param("id")
	var req struct {
		AdminID uint `json:"admin_id"`
	}
	c.ShouldBindJSON(&req)

	var room models.ChatRoom
	if err := config.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	room.AdminID = &req.AdminID
	room.Status = "ASSIGNED" // Or keep OPEN
	config.DB.Save(&room)

	// Notify User via System Message
	sysMsg := models.Message{
		ChatRoomID: room.ID,
		SenderID:   0,
		SenderType: "SYSTEM",
		Content:    "Admin telah bergabung ke percakapan.",
		CreatedAt:  time.Now(),
	}
	config.DB.Create(&sysMsg)

	c.JSON(http.StatusOK, gin.H{"message": "Joined chat", "data": room})
}

// ResolveChat (Close)
func ResolveChat(c *gin.Context) {
	roomID := c.Param("id")
	config.DB.Model(&models.ChatRoom{}).Where("id = ?", roomID).Update("status", "RESOLVED")
	c.JSON(http.StatusOK, gin.H{"message": "Chat resolved"})
}

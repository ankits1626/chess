package websocket

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 1 * 1024 * 1024 // 1 MB
)

// Client represents a WebSocket client connection.
type Client struct {
	// Unique client ID
	ID string

	// User ID from authentication
	UserID string

	// Current game ID (if in a game)
	GameID string

	// WebSocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Reference to hub
	hub *Hub

	// Current room (if joined)
	room *Room

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewClient creates a new WebSocket client.
func NewClient(conn *websocket.Conn, hub *Hub, userID string) *Client {
	return &Client{
		ID:     generateID(),
		UserID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    hub,
	}
}

// readPump pumps messages from the WebSocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(ctx context.Context, messageHandler func(context.Context, *Client, []byte)) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.ID, err)
			}
			break
		}

		// Handle message with context
		messageHandler(ctx, c, message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Write error for client %s: %v", c.ID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to the client.
func (c *Client) SendMessage(msg *Message) error {
	data, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("client send buffer full")
	}
}

// JoinRoom adds client to a room.
// Returns error if client is already in a room or room is full.
func (c *Client) JoinRoom(room *Room) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already in a room
	if c.room != nil {
		return fmt.Errorf("client already in room %s", c.room.ID)
	}

	// Check room capacity (chess is max 2 players)
	if room.ClientCount() >= 2 {
		return fmt.Errorf("room is full (max 2 players)")
	}

	c.room = room
	c.GameID = room.ID
	room.AddClient(c)

	log.Printf("Client %s (user %s) joined room %s", c.ID, c.UserID, room.ID)
	return nil
}

// LeaveRoom removes client from current room.
func (c *Client) LeaveRoom() {
	c.mu.Lock()
	room := c.room
	c.room = nil
	c.GameID = ""
	c.mu.Unlock()

	if room != nil {
		room.RemoveClient(c)
		log.Printf("Client %s (user %s) left room %s", c.ID, c.UserID, room.ID)
	}
}

// GetRoom returns the current room (thread-safe).
func (c *Client) GetRoom() *Room {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.room
}

// SetGameID sets the game ID for this client (thread-safe).
func (c *Client) SetGameID(gameID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.GameID = gameID
}

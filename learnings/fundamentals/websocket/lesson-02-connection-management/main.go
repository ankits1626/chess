package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10 // 54 seconds

	// Maximum message size allowed from peer
	maxMessageSize = 1 * 1024 * 1024 // 1 MB (perfect for chess-coach)

	// Buffer sizes
	readBufferSize  = 4096
	writeBufferSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBufferSize,
	WriteBufferSize: writeBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for learning
	},
}

// Client represents a WebSocket client connection
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	registry *ConnectionRegistry
	id       string
}

// ConnectionRegistry tracks all active connections
type ConnectionRegistry struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

// NewConnectionRegistry creates a new connection registry
func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

// Run starts the connection registry event loop
func (r *ConnectionRegistry) Run() {
	for {
		select {
		case client := <-r.register:
			r.mu.Lock()
			r.clients[client] = true
			r.mu.Unlock()
			log.Printf("✅ Client %s registered. Total connections: %d", client.id, r.Count())

		case client := <-r.unregister:
			r.mu.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			r.mu.Unlock()
			log.Printf("👋 Client %s unregistered. Total connections: %d", client.id, r.Count())

		case message := <-r.broadcast:
			r.mu.RLock()
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					// Client's send buffer is full, close it
					close(client.send)
					delete(r.clients, client)
				}
			}
			r.mu.RUnlock()
		}
	}
}

// Count returns the number of active connections
func (r *ConnectionRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

// Shutdown gracefully closes all connections
func (r *ConnectionRegistry) Shutdown() {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Printf("🛑 Shutting down... Closing %d connections", len(r.clients))

	for client := range r.clients {
		// Send close message with reason
		message := websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down")
		client.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(writeWait))
		client.conn.Close()
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.registry.unregister <- c
		c.conn.Close()
	}()

	// Configure connection
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// Set pong handler - when pong received, reset read deadline
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Printf("🏓 Received PONG from client %s", c.id)
		return nil
	})

	// Read messages from client
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ Client %s: unexpected close error: %v", c.id, err)
			} else {
				log.Printf("👋 Client %s disconnected normally", c.id)
			}
			break
		}

		// Log received message
		if messageType == websocket.TextMessage {
			log.Printf("📨 Client %s sent TEXT: %s", c.id, string(message))
		} else {
			log.Printf("📨 Client %s sent BINARY: %d bytes", c.id, len(message))
		}

		// Echo message back to the sender
		response := fmt.Sprintf("Server echo: %s (at %s)", string(message), time.Now().Format("15:04:05"))
		select {
		case c.send <- []byte(response):
		default:
			log.Printf("⚠️ Client %s send buffer full", c.id)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	// Create ticker for pings
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
				// Channel closed, send close message
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send message to client
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("❌ Error writing to client %s: %v", c.id, err)
				return
			}
			log.Printf("📤 Sent to client %s: %s", c.id, string(message))

		case <-ticker.C:
			// Send ping to client
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("❌ Failed to send PING to client %s: %v", c.id, err)
				return
			}
			log.Printf("🏓 Sent PING to client %s", c.id)
		}
	}
}

// handleWebSocket handles WebSocket upgrade and connection
func handleWebSocket(registry *ConnectionRegistry, w http.ResponseWriter, r *http.Request) {
	// Generate client ID
	clientID := fmt.Sprintf("%s-%d", r.RemoteAddr, time.Now().Unix())

	log.Printf("📞 New connection request from %s (ID: %s)", r.RemoteAddr, clientID)

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade connection: %v", err)
		return
	}

	// Create client
	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		registry: registry,
		id:       clientID,
	}

	// Register client
	registry.register <- client

	// Send welcome message
	welcomeMsg := fmt.Sprintf("Welcome! You are client %s. Connected at %s", clientID, time.Now().Format("15:04:05"))
	client.send <- []byte(welcomeMsg)

	// Start goroutines for this client
	go client.writePump() // Must start writePump first (it sends pings)
	go client.readPump()  // Then start readPump (it handles pongs)
}

// handleHome serves the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Lesson 2 - Connection Management</title>
</head>
<body>
    <h1>WebSocket Lesson 2: Connection Management</h1>
    <p>Server is running with production-ready connection management! 🚀</p>
    <h2>Features Implemented:</h2>
    <ul>
        <li>✅ Ping/Pong heartbeats (every 54 seconds)</li>
        <li>✅ Read/Write deadlines (60s pong wait, 10s write wait)</li>
        <li>✅ Connection registry tracking</li>
        <li>✅ Graceful shutdown</li>
        <li>✅ Proper error handling</li>
    </ul>
    <h2>Connection Info:</h2>
    <ul>
        <li><strong>WebSocket URL:</strong> ws://localhost:8080/ws</li>
        <li><strong>Test Client:</strong> <a href="/client.html">Open client.html</a></li>
    </ul>
    <h2>Watch the logs!</h2>
    <p>Check your terminal to see:</p>
    <ul>
        <li>🏓 PING/PONG messages every 54 seconds</li>
        <li>📨 Message echoes</li>
        <li>✅ Connection/disconnection events</li>
        <li>📊 Active connection count</li>
    </ul>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func main() {
	// Print startup banner
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║   Lesson 2: Connection Management             ║")
	fmt.Println("║   Production-Ready WebSocket Server           ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("⚙️  Configuration:\n")
	fmt.Printf("   • Ping interval:    %v\n", pingPeriod)
	fmt.Printf("   • Pong timeout:     %v\n", pongWait)
	fmt.Printf("   • Write timeout:    %v\n", writeWait)
	fmt.Printf("   • Max message size: %d MB\n", maxMessageSize/(1024*1024))
	fmt.Println()

	// Create connection registry
	registry := NewConnectionRegistry()
	go registry.Run()

	// Register HTTP handlers
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(registry, w, r)
	})

	// Start server
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println()
		log.Println("🛑 Shutdown signal received...")

		// Gracefully close all WebSocket connections
		registry.Shutdown()

		// Shutdown HTTP server
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}

		log.Println("✅ Server stopped gracefully")
		os.Exit(0)
	}()

	// Start HTTP server
	addr := ":8080"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", addr)
	fmt.Println()
	fmt.Println("📋 Endpoints:")
	fmt.Println("   • Home page:     http://localhost:8080/")
	fmt.Println("   • WebSocket:     ws://localhost:8080/ws")
	fmt.Println("   • Test client:   http://localhost:8080/client.html")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("   • Watch for 🏓 PING/PONG messages every 54 seconds")
	fmt.Println("   • Try closing your browser to see timeout detection")
	fmt.Println("   • Press Ctrl+C to see graceful shutdown")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}
}

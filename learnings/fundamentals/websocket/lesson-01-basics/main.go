package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed (10 MB)
	// Adjust based on your use case:
	// - Chat: 64 KB
	// - API data: 1 MB
	// - Images: 5-10 MB
	maxMessageSize = 10 * 1024 * 1024 // 10 MB

	// Buffer sizes for reading/writing
	readBufferSize  = 4096 // 4 KB
	writeBufferSize = 4096 // 4 KB
)

// upgrader is configured to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBufferSize,
	WriteBufferSize: writeBufferSize,
	// CheckOrigin allows all connections for learning purposes
	// In production, implement proper CORS checking!
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handleWebSocket handles WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Log the connection attempt
	log.Printf("📞 New connection request from %s", r.RemoteAddr)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Set maximum message size to prevent memory exhaustion attacks
	conn.SetReadLimit(maxMessageSize)

	log.Printf("✅ WebSocket connection established with %s (max message size: %d MB)", r.RemoteAddr, maxMessageSize/(1024*1024))

	// Send welcome message
	welcomeMsg := fmt.Sprintf("Welcome! Connected at %s", time.Now().Format("15:04:05"))
	err = conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))
	if err != nil {
		log.Printf("❌ Error sending welcome message: %v", err)
		return
	}

	// Listen for messages from client
	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Check if it's a normal closure
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ Error reading message: %v", err)
			} else {
				log.Printf("👋 Client disconnected normally")
			}
			break
		}

		// Handle different message types
		if messageType == websocket.BinaryMessage {
			// Binary message (images, files, etc.)
			log.Printf("📨 Received BINARY message: %d bytes", len(message))

			// Echo binary data back as-is
			err = conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Printf("❌ Error sending binary message: %v", err)
				break
			}

			log.Printf("📤 Sent BINARY echo: %d bytes", len(message))
		} else {
			// Text message
			log.Printf("📨 Received (%s): %s", getMessageTypeName(messageType), string(message))

			// Echo the message back to client
			responseMsg := fmt.Sprintf("Echo: %s (received at %s)", string(message), time.Now().Format("15:04:05"))
			err = conn.WriteMessage(websocket.TextMessage, []byte(responseMsg))
			if err != nil {
				log.Printf("❌ Error sending message: %v", err)
				break
			}

			log.Printf("📤 Sent: %s", responseMsg)
		}
	}

	log.Printf("🔌 Connection closed with %s", r.RemoteAddr)
}

// getMessageTypeName returns human-readable message type name
func getMessageTypeName(messageType int) string {
	switch messageType {
	case websocket.TextMessage:
		return "TEXT"
	case websocket.BinaryMessage:
		return "BINARY"
	case websocket.CloseMessage:
		return "CLOSE"
	case websocket.PingMessage:
		return "PING"
	case websocket.PongMessage:
		return "PONG"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", messageType)
	}
}

// handleHome serves a simple home page with connection info
func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Lesson 1 - WebSocket Basics</title>
</head>
<body>
    <h1>WebSocket Lesson 1: Basics</h1>
    <p>Server is running! 🚀</p>
    <h2>Connection Info:</h2>
    <ul>
        <li><strong>WebSocket URL:</strong> ws://localhost:8080/ws</li>
        <li><strong>Test Client:</strong> <a href="/client.html">Open client.html</a></li>
    </ul>
    <h2>What this server does:</h2>
    <ol>
        <li>Accepts WebSocket connections at /ws</li>
        <li>Sends a welcome message</li>
        <li>Echoes back any message you send</li>
        <li>Logs everything to console</li>
    </ol>
    <h2>Try it:</h2>
    <pre>
// JavaScript in browser console:
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => console.log('Connected!');
ws.onmessage = (e) => console.log('Received:', e.data);
ws.send('Hello WebSocket!');
    </pre>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func main() {
	// Print startup banner
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   Lesson 1: WebSocket Basics           ║")
	fmt.Println("║   Simple Echo Server                   ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// Register HTTP handlers
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)

	// Serve static files (for client.html)
	http.Handle("/client.html", http.FileServer(http.Dir(".")))

	// Start server
	addr := ":8080"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", addr)
	fmt.Println()
	fmt.Println("📋 Endpoints:")
	fmt.Println("   • Home page:     http://localhost:8080/")
	fmt.Println("   • WebSocket:     ws://localhost:8080/ws")
	fmt.Println("   • Test client:   http://localhost:8080/client.html")
	fmt.Println()
	fmt.Println("💡 Tip: Watch this terminal for connection logs!")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(addr, nil))
}

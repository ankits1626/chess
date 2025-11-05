package main

import (
	"fmt"
	"time"
)

// Simplified Hub (like websocket hub)
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// Simplified Client
type Client struct {
	id   string
	hub  *Hub
	send chan []byte
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 10),
	}
}

func (h *Hub) run() {
	fmt.Println("🚀 Hub started\n")
	for {
		select {
		case client := <-h.register:
			// ← Receive client from register channel
			h.clients[client] = true
			fmt.Printf("📥 Hub: Registered client '%s' (total: %d)\n", client.id, len(h.clients))

		case client := <-h.unregister:
			// ← Receive client from unregister channel
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Printf("📤 Hub: Unregistered client '%s' (total: %d)\n", client.id, len(h.clients))
			}

		case message := <-h.broadcast:
			// ← Receive message from broadcast channel
			fmt.Printf("\n📢 Hub: Broadcasting message to %d clients\n", len(h.clients))
			for client := range h.clients {
				select {
				case client.send <- message:
					// → Send message to client's send channel
					fmt.Printf("   ✓ Sent to '%s'\n", client.id)
				default:
					// Client buffer full, skip
					fmt.Printf("   ✗ Client '%s' buffer full\n", client.id)
				}
			}
			fmt.Println()
		}
	}
}

func (c *Client) writePump() {
	fmt.Printf("  [%s] writePump started\n", c.id)
	for {
		message, ok := <-c.send // ← Receive from send channel
		if !ok {
			fmt.Printf("  [%s] writePump: channel closed, exiting\n", c.id)
			return
		}
		fmt.Printf("  [%s] 📨 Received: %s\n", c.id, string(message))
	}
}

func main() {
	fmt.Println("=== Example 6: Websocket Hub Simulation ===\n")

	// Create hub
	hub := newHub()
	go hub.run()

	// Create clients
	client1 := &Client{id: "alice", hub: hub, send: make(chan []byte, 10)}
	client2 := &Client{id: "bob", hub: hub, send: make(chan []byte, 10)}
	client3 := &Client{id: "charlie", hub: hub, send: make(chan []byte, 10)}

	// Register clients
	hub.register <- client1 // → Send to register channel
	hub.register <- client2
	hub.register <- client3
	time.Sleep(100 * time.Millisecond)

	// Start write pumps
	go client1.writePump()
	go client2.writePump()
	go client3.writePump()
	time.Sleep(100 * time.Millisecond)

	// Broadcast messages
	fmt.Println("--- Broadcasting Messages ---\n")
	hub.broadcast <- []byte("Game started!")
	time.Sleep(200 * time.Millisecond)

	hub.broadcast <- []byte("Move: e2e4")
	time.Sleep(200 * time.Millisecond)

	// Unregister one client
	fmt.Println("--- Client Disconnect ---\n")
	hub.unregister <- client2
	time.Sleep(200 * time.Millisecond)

	// Broadcast again (only 2 clients now)
	hub.broadcast <- []byte("Move: e7e5")
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   Hub uses channels to coordinate all clients")
	fmt.Println("   register/unregister = control channels")
	fmt.Println("   broadcast = message channel")
	fmt.Println("   client.send = per-client message queue")
	fmt.Println("\n   This is EXACTLY how your websocket hub works!")
}

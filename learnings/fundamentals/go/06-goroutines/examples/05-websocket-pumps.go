package main

import (
	"fmt"
	"time"
)

// Simulated Client
type Client struct {
	id   string
	send chan []byte
}

// readPump simulates reading from websocket
func (c *Client) readPump() {
	fmt.Printf("[%s] 📖 readPump started\n", c.id)

	// Simulate reading 3 messages from client
	messages := []string{"move: e2e4", "move: e7e5", "move: Nf3"}

	for _, msg := range messages {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("[%s] 📨 readPump: Received from websocket: %s\n", c.id, msg)

		// In real code: would send to hub
		// hub.broadcast <- []byte(msg)
	}

	fmt.Printf("[%s] 📖 readPump: Client disconnected\n", c.id)
}

// writePump simulates writing to websocket
func (c *Client) writePump() {
	fmt.Printf("[%s] ✍️  writePump started\n", c.id)

	// Continuously receive from send channel
	for message := range c.send {
		fmt.Printf("[%s] 📤 writePump: Sending to websocket: %s\n", c.id, string(message))
		time.Sleep(200 * time.Millisecond) // Simulate network delay
	}

	fmt.Printf("[%s] ✍️  writePump: Send channel closed, exiting\n", c.id)
}

func main() {
	fmt.Println("=== Example 5: Websocket Read/Write Pumps ===\n")

	// Create client
	client := &Client{
		id:   "player1",
		send: make(chan []byte, 10),
	}

	fmt.Println("🔌 Client connected, starting pumps...\n")

	// Start TWO goroutines for this client
	go client.readPump()  // Goroutine 1: Reads from websocket
	go client.writePump() // Goroutine 2: Writes to websocket

	// Simulate sending messages to client from hub
	time.Sleep(300 * time.Millisecond)
	client.send <- []byte("Welcome to the game!")

	time.Sleep(600 * time.Millisecond)
	client.send <- []byte("Opponent joined")

	time.Sleep(700 * time.Millisecond)
	client.send <- []byte("Your turn")

	// Wait for readPump to finish
	time.Sleep(1 * time.Second)

	// Close send channel (client disconnect)
	close(client.send)
	time.Sleep(300 * time.Millisecond)

	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   readPump + writePump run CONCURRENTLY (same time)")
	fmt.Println("   readPump:  Blocks waiting for websocket messages")
	fmt.Println("   writePump: Blocks waiting for messages to send")
	fmt.Println("   They MUST be in separate goroutines!")
	fmt.Println("\n   This is WHY you see:")
	fmt.Println("   go client.readPump()")
	fmt.Println("   go client.writePump()")
}

package main

import (
	"fmt"
	"time"
)

// Hub manages all clients
type Hub struct {
	clients map[string]bool
}

func (h *Hub) Run() {
	fmt.Println("🚀 Hub.Run() started (infinite loop)")

	// In real code: this runs forever
	for i := 1; i <= 3; i++ {
		fmt.Printf("   Hub: Loop iteration %d - waiting for events...\n", i)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("   Hub: (In real code, this never exits)")
}

func startHTTPServer() {
	fmt.Println("🌐 HTTP Server started (infinite loop)")

	// In real code: http.ListenAndServe blocks forever
	for i := 1; i <= 3; i++ {
		fmt.Printf("   Server: Handling request %d...\n", i)
		time.Sleep(600 * time.Millisecond)
	}

	fmt.Println("   Server: (In real code, this never exits)")
}

func main() {
	fmt.Println("=== Example 6: Hub + Server Pattern ===\n")

	// Create hub
	hub := &Hub{
		clients: make(map[string]bool),
	}

	// --- Problem: Both block forever! ---
	fmt.Println("--- THE PROBLEM (commented out) ---")
	fmt.Println("If we run sequentially:")
	fmt.Println("  hub.Run()              // ← Blocks forever!")
	fmt.Println("  startHTTPServer()      // ← Never reached!\n")

	// --- Solution: Run hub in goroutine ---
	fmt.Println("--- THE SOLUTION ---")
	fmt.Println("Run hub in background:\n")

	go hub.Run() // ← Start hub in goroutine (background)

	// Now main can continue to start HTTP server
	time.Sleep(100 * time.Millisecond) // Let hub start

	fmt.Println("Main: Hub is running in background, starting server...\n")

	// In real code: this would be http.ListenAndServe (also blocks)
	startHTTPServer()

	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   Hub.Run() has infinite loop → must run in goroutine")
	fmt.Println("   This frees main to start HTTP server")
	fmt.Println("\n   This is WHY you see:")
	fmt.Println("   go hub.Run()")
	fmt.Println("   http.ListenAndServe(...)")
}

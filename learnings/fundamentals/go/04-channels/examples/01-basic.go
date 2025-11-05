package main

import "fmt"

func main() {
	fmt.Println("=== Example 1: Basic Channel Send & Receive ===\n")

	// Create a channel for strings
	messages := make(chan string)

	// Send in a goroutine (background thread)
	go func() {
		fmt.Println("📤 Goroutine: Sending 'ping'...")
		messages <- "ping" // ← Send INTO channel
		messages <- "ping 2"
		fmt.Println("📤 Goroutine: Sent!")
	}()

	// Receive in main goroutine
	fmt.Println("📥 Main: Waiting to receive...")
	msg := <-messages // ← Receive FROM channel
	fmt.Println("📥 Main: Received:", msg)

	fmt.Println("\n--- Example 2: Multiple Messages ---\n")

	// Send multiple messages
	go func() {
		messages <- "first"
		messages <- "second"
		messages <- "third"
	}()

	// Receive multiple messages
	fmt.Println("Message 1:", <-messages)
	fmt.Println("Message 2:", <-messages)
	fmt.Println("Message 3:", <-messages)

	fmt.Println("\n✅ Key takeaway: <- direction shows data flow!")
	fmt.Println("   ch <- v    = Send v INTO channel")
	fmt.Println("   v := <-ch  = Receive FROM channel into v")
}

package main

import "fmt"

func main() {
	fmt.Println("=== Example 2: Buffered vs Unbuffered Channels ===\n")

	// --- Unbuffered Channel ---
	fmt.Println("--- Unbuffered Channel ---")
	unbuffered := make(chan string)

	// This would DEADLOCK if not in goroutine!
	// unbuffered <- "hello"  // ❌ Blocks forever - no receiver

	go func() {
		unbuffered <- "hello" // ✅ OK - someone will receive
	}()

	msg := <-unbuffered
	fmt.Println("Received from unbuffered:", msg)

	// --- Buffered Channel ---
	fmt.Println("\n--- Buffered Channel (size 3) ---")
	buffered := make(chan string, 3)

	// Can send without goroutine (has space in buffer)
	buffered <- "first"  // ✅ OK - buffer has space
	buffered <- "second" // ✅ OK - buffer has space
	buffered <- "third"  // ✅ OK - buffer has space
	fmt.Println("Sent 3 messages to buffer")

	// buffered <- "fourth"  // ❌ Would BLOCK - buffer full!

	// Receive from buffer
	fmt.Println("Received:", <-buffered) // "first"
	fmt.Println("Received:", <-buffered) // "second"
	fmt.Println("Received:", <-buffered) // "third"

	// --- Chess Coach Pattern ---
	fmt.Println("\n--- Websocket Pattern (Buffered Send Queue) ---")

	// Like client.send in websocket code
	sendQueue := make(chan []byte, 256) // Buffer 256 messages

	// Can queue messages without blocking
	sendQueue <- []byte("move: e2e4")
	sendQueue <- []byte("move: e7e5")
	sendQueue <- []byte("move: Nf3")
	fmt.Println("Queued 3 moves")

	// Process queue
	fmt.Println("Processing:", string(<-sendQueue))
	fmt.Println("Processing:", string(<-sendQueue))
	fmt.Println("Processing:", string(<-sendQueue))

	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   Unbuffered = synchronous (sender waits for receiver)")
	fmt.Println("   Buffered   = asynchronous (sender doesn't wait if space available)")
}

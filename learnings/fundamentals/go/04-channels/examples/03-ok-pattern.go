package main

import "fmt"

func main() {
	fmt.Println("=== Example 3: The 'ok' Pattern (Channel Closed Detection) ===\n")

	messages := make(chan string, 3)

	// Send some messages and close
	go func() {
		messages <- "hello"
		messages <- "world"
		messages <- "goodbye"
		close(messages) // Signal: no more data
		fmt.Println("📤 Sender: Closed channel")
	}()

	// --- Method 1: Check with ok ---
	fmt.Println("--- Using ok to detect closure ---")
	for i := 1; i <= 4; i++ {
		msg, ok := <-messages
		if !ok {
			fmt.Printf("Attempt %d: Channel is closed\n", i)
			break
		}
		fmt.Printf("Attempt %d: Received '%s'\n", i, msg)
	}

	// --- Method 2: Range (automatic) ---
	fmt.Println("\n--- Using range (cleaner) ---")

	messages2 := make(chan string)
	go func() {
		messages2 <- "first"
		messages2 <- "second"
		messages2 <- "third"
		close(messages2)
	}()

	// range automatically stops when channel closes
	for msg := range messages2 {
		fmt.Println("Received:", msg)
	}
	fmt.Println("Channel closed, loop exited")

	// --- Websocket Pattern ---
	fmt.Println("\n--- Websocket writePump Pattern ---")

	clientSend := make(chan []byte, 2)

	go func() {
		clientSend <- []byte("message 1")
		clientSend <- []byte("message 2")
		close(clientSend) // Hub closes when client disconnects
	}()

	// writePump pattern
	for {
		message, ok := <-clientSend
		if !ok {
			fmt.Println("Client send channel closed - disconnecting")
			break
		}
		fmt.Printf("Sending to websocket: %s\n", string(message))
	}

	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   msg, ok := <-ch")
	fmt.Println("   ok = true  → channel open, msg is valid")
	fmt.Println("   ok = false → channel closed, no more data")
}

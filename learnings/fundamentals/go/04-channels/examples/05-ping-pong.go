package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Example 5: Ping-Pong Game (Two-Way Communication) ===\n")

	ping := make(chan string)
	pong := make(chan string)

	// Player 1 (receives ping, sends pong)
	go func() {
		for i := 0; i < 3; i++ {
			// Receive from ping channel
			msg := <-ping
			fmt.Printf("  Player 1 received: %s\n", msg)
			time.Sleep(200 * time.Millisecond)

			// Send to pong channel
			pong <- "pong"
			fmt.Printf("  Player 1 sent: pong\n\n")
		}
	}()

	// Player 2 (receives pong, sends ping)
	go func() {
		for i := 0; i < 3; i++ {
			// Receive from pong channel
			msg := <-pong
			fmt.Printf("    Player 2 received: %s\n", msg)
			time.Sleep(200 * time.Millisecond)

			// Send to ping channel (if not last round)
			if i < 2 {
				ping <- "ping"
				fmt.Printf("    Player 2 sent: ping\n\n")
			}
		}
	}()

	// Start the game!
	fmt.Println("🏓 Starting game...\n")
	ping <- "ping"
	fmt.Println("Game started with 'ping'\n")

	// Wait for game to complete
	time.Sleep(2 * time.Second)

	fmt.Println("🏁 Game over!\n")

	fmt.Println("✅ Key takeaway:")
	fmt.Println("   Channels enable back-and-forth communication")
	fmt.Println("   Each goroutine blocks waiting for messages")
	fmt.Println("   This creates synchronized turn-taking")
}

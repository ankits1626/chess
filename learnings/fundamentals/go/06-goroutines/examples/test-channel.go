package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	// Start goroutine that waits for message
	go func() {
		fmt.Println("Goroutine: Waiting for message...")

		for msg := range ch { // ← BLOCKS here waiting

			fmt.Println("Goroutine: Got message:", msg)
		}
		fmt.Println("Goroutine: Channel closed, exiting")
	}()

	// Main goroutine
	fmt.Println("Main: Sleeping for 2 seconds...")
	time.Sleep(2 * time.Second)

	fmt.Println("Main: Sending message now!")
	ch <- "Hello!" // ← This unblocks the goroutine above

	close(ch)

	fmt.Println("Main: Sleeping for 2 seconds again...")
	time.Sleep(2 * time.Second)

	fmt.Println("Main: Sending message now!")
	ch <- "World!"
	time.Sleep(100 * time.Millisecond)
	// close(ch)
	time.Sleep(100 * time.Millisecond) // Let it print
	fmt.Println("Main: Done!")
}

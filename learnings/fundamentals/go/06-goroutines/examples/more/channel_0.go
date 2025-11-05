package main

import "fmt"

func main() {
	// Create a channel
	ch := make(chan string)

	// Send and receive
	go func() {
		ch <- "Hello from goroutine!" // Send to channel
	}()

	msg := <-ch // Receive from channel
	fmt.Println(msg)
}

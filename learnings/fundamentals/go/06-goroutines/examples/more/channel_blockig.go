package main

import (
	"fmt"
	"time"
)

func main() {
	n := make(chan int)

	go func() {
		n <- 1
		println("N is receibed by the sender")
	}()

	time.Sleep(2 * time.Second)
	received := <-n
	fmt.Printf("Received %d from goroutine", received)
}

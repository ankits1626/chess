package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Started")
	ch := make(chan int, 2)

	ch <- 1
	fmt.Println("Sent 1")
	ch <- 2
	fmt.Println("Sent 2")

	// ch <- 2
	// fmt.Println("Sent 3")

	go func() {
		fmt.Println("Trying to send 3...")
		ch <- 3 // BLOCKS here until someone receives
		fmt.Println("Sent 3!")
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("Main sleeping... goroutine is blocked!")

}

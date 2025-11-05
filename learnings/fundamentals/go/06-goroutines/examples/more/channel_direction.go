package main

import (
	"fmt"
	"time"
)

func send(ch chan<- int) {
	ch <- 1
}

func receive(ch <-chan int) {
	msg := <-ch
	fmt.Println("Received message: %d", msg)
}

func main() {
	fmt.Println("Start: ", time.Now())
	ch := make(chan int)

	go send(ch)
	receive(ch)
	time.Sleep(1)

}

package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(1)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(2)
		ch2 <- 2
	}()

	select {
	case msh := <-ch1:
		fmt.Printf("Msg = %d", msh)
	case msh2 := <-ch2:
		fmt.Printf("Msg = %d", msh2)
	}

}

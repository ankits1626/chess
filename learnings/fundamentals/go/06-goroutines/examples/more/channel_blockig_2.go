package main

import "fmt"

func main() {
	n := make(chan int)

	n <- 1

	fmt.Printf("This line never runs")
}

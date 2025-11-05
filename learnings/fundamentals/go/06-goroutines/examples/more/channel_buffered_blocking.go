package main

import "fmt"

func main() {
	n := make(chan int, 2)
	n <- 2
	n <- 3
	fmt.Printf("This prints\n")
	n <- 4
	fmt.Printf("this does not print")
}

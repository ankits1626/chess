package main

import (
	"fmt"
	"time"
)

func count(name string) {
	for i := 1; i <= 5; i++ {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// WITHOUT goroutines (sequential)
	fmt.Println("=== Sequential ===")
	count("First")
	count("Second")

	fmt.Println("\n=== Concurrent ===")
	// WITH goroutines (concurrent)
	go count("Goroutine-1")
	go count("Goroutine-2")

	// Wait for goroutines to finish
	time.Sleep(600 * time.Millisecond)
}

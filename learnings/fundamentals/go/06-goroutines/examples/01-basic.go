package main

import (
	"fmt"
	"time"
)

func sayHello(name string) {
	fmt.Printf("Hello from %s!\n", name)
}

func countTo5(name string) {
	for i := 1; i <= 5; i++ {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("=== Example 1: Basic Goroutine ===\n")

	// WITHOUT goroutine - sequential
	fmt.Println("--- Sequential (no 'go') ---")
	sayHello("Alice")
	sayHello("Bob")
	fmt.Println()

	// WITH goroutine - concurrent
	fmt.Println("--- Concurrent (with 'go') ---")
	go sayHello("Charlie") // Starts in background
	go sayHello("Diana")   // Starts in background

	time.Sleep(100 * time.Millisecond) // Wait for goroutines to finish
	fmt.Println()

	// Multiple goroutines running at same time
	fmt.Println("--- Two Counters Running Simultaneously ---")
	go countTo5("Worker 1")
	go countTo5("Worker 2")

	time.Sleep(600 * time.Millisecond) // Wait for both to finish
	fmt.Println()

	fmt.Println("✅ Key takeaway:")
	fmt.Println("   'go' keyword starts a function in background (new goroutine)")
	fmt.Println("   Multiple goroutines run concurrently (at the same time)")
}

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Example 2: Anonymous Goroutines ===\n")

	// --- Pattern 1: Simple anonymous function ---
	fmt.Println("--- Anonymous Function as Goroutine ---")
	go func() {
		fmt.Println("  Running in goroutine!")
	}() // ← Note the () to call it immediately

	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	// --- Pattern 2: With parameters ---
	fmt.Println("--- With Parameters ---")
	message := "Hello from parameter"
	go func(msg string) {
		fmt.Printf("  Goroutine received: %s\n", msg)
	}(message) // ← Pass parameter here

	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	// --- Pattern 3: Closure (accesses outer variables) ---
	fmt.Println("--- Closure (Accessing Outer Variables) ---")
	name := "Alice"
	score := 100

	go func() {
		fmt.Printf("  Player: %s, Score: %d\n", name, score)
	}() // Accesses 'name' and 'score' from outer scope

	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	// --- Pattern 4: Loop with goroutines (CAREFUL!) ---
	fmt.Println("--- Loop with Goroutines (WRONG WAY) ---")
	for i := 0; i < 3; i++ {
		go func() {
			fmt.Printf("  ❌ i = %d (might all print 3!)\n", i)
		}()
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	fmt.Println("--- Loop with Goroutines (RIGHT WAY) ---")
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Printf("  ✅ id = %d (each gets its own copy)\n", id)
		}(i) // Pass i as parameter
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	// --- Websocket Pattern ---
	fmt.Println("--- Websocket Server Pattern ---")
	go func() {
		fmt.Println("  🚀 Starting HTTP server in background...")
		// In real code: http.ListenAndServe(":8080", nil)
		time.Sleep(200 * time.Millisecond)
		fmt.Println("  ✅ Server running")
	}()

	fmt.Println("  Main goroutine continues...")
	time.Sleep(300 * time.Millisecond)
	fmt.Println()

	fmt.Println("✅ Key takeaway:")
	fmt.Println("   go func() {...}()  starts anonymous function as goroutine")
	fmt.Println("   Pass loop variables as parameters to avoid bugs!")
}

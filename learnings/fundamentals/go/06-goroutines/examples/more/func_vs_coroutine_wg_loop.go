package main

import (
	"fmt"
	"sync"
	"time"
)

func count(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 5; i++ {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("=== Concurrent with WaitGroup and Loop ===")

	numGoroutines := 5

	wg.Add(numGoroutines) // Add all goroutines upfront

	for i := 1; i <= numGoroutines; i++ {
		goroutineName := fmt.Sprintf("Goroutine-%d", i)
		go count(goroutineName, &wg)
	}

	wg.Wait() // Wait for all goroutines to finish

	fmt.Println("\nAll goroutines finished!")
}

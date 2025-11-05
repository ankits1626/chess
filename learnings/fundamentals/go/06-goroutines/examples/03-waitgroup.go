package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when function returns

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Duration(id*100) * time.Millisecond) // Simulate work
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	fmt.Println("=== Example 3: WaitGroup (Proper Synchronization) ===\n")

	// --- Problem: time.Sleep is unreliable ---
	fmt.Println("--- WITHOUT WaitGroup (Unreliable) ---")
	for i := 1; i <= 10; i++ {
		go func(id int) {
			// fmt.Printf("  Task %d running\n", id)
			// time.Sleep(0.001 * time.Millisecond)
			fmt.Printf("Task %d done\n", id)
		}(i)
	}
	time.Sleep(150 * time.Millisecond) // Hope this is enough!
	fmt.Println("  Main: Assuming all done (but not sure!)\n")

	// --- Solution: WaitGroup ---
	fmt.Println("--- WITH WaitGroup (Reliable) ---")

	var wg sync.WaitGroup

	// Start 5 workers
	for i := 1; i <= 5; i++ {
		wg.Add(1) // Register 1 goroutine
		go worker(i, &wg)
	}

	fmt.Println("Main: Waiting for all workers...\n")
	wg.Wait() // Block until all goroutines call wg.Done()
	fmt.Println("\nMain: All workers confirmed done!\n")

	// --- Real-world pattern ---
	fmt.Println("--- Real-World Pattern: Processing Items ---")

	items := []string{"task1", "task2", "task3", "task4", "task5"}
	var wg2 sync.WaitGroup

	for _, item := range items {
		wg2.Add(1)
		go func(task string) {
			defer wg2.Done()
			fmt.Printf("  Processing: %s\n", task)
			time.Sleep(100 * time.Millisecond)
		}(item) // Pass item as parameter!
	}

	wg2.Wait()
	fmt.Println("  All items processed!\n")

	fmt.Println("✅ Key takeaway:")
	fmt.Println("   WaitGroup: Proper way to wait for multiple goroutines")
	fmt.Println("   wg.Add(1)   → Register a goroutine")
	fmt.Println("   wg.Done()   → Mark goroutine complete")
	fmt.Println("   wg.Wait()   → Block until all done")
}

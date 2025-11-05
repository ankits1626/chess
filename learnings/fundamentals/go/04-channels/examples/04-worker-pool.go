package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d: Processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- job * 2                  // Send result
	}
	fmt.Printf("Worker %d: Finished (jobs channel closed)\n", id)
}

func main() {
	fmt.Println("=== Example 4: Worker Pool Pattern ===\n")

	// Create channels
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 9 jobs
	fmt.Println("📤 Sending 9 jobs...\n")
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs) // No more jobs

	// Collect results
	fmt.Println("\n📥 Collecting results...")
	for r := 1; r <= 9; r++ {
		result := <-results
		fmt.Printf("Result %d: %d\n", r, result)
	}

	time.Sleep(200 * time.Millisecond) // Let workers finish printing

	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   Multiple workers receive from SAME jobs channel")
	fmt.Println("   Each job goes to only ONE worker (distributed)")
	fmt.Println("   This is how you parallelize work in Go!")
}

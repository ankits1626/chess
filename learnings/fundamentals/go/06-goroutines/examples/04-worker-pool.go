package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d: Processing job %d\n", id, job)
		time.Sleep(500 * time.Millisecond) // Simulate work
		results <- job * 2                  // Send result
	}
	fmt.Printf("Worker %d: No more jobs, exiting\n", id)
}

func main() {
	fmt.Println("=== Example 4: Worker Pool Pattern ===\n")

	// Create channels
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Start 3 workers
	var wg sync.WaitGroup
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send 9 jobs
	fmt.Println("📤 Sending 9 jobs to worker pool...\n")
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs) // No more jobs - workers will exit

	// Collect results in separate goroutine
	go func() {
		wg.Wait()      // Wait for all workers to finish
		close(results) // Then close results channel
	}()

	// Receive all results
	fmt.Println("\n📥 Collecting results...")
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println("\n✅ All work complete!")
	fmt.Println("\n✅ Key takeaway:")
	fmt.Println("   Worker pool: Fixed number of goroutines processing jobs")
	fmt.Println("   Prevents creating too many goroutines")
	fmt.Println("   Common pattern for parallel processing")
}

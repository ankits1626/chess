package main

import (
	"fmt"
	"sync"
	"time"
)

// Exercise 1: Start a simple goroutine
func exercise1() {
	fmt.Println("\n=== Exercise 1: Start a Goroutine ===")

	// TODO: Create a function that prints "Hello from goroutine!"
	// TODO: Call it with 'go' keyword
	// TODO: Wait 100ms for it to finish

	go func() {
		fmt.Println("Hello from goroutine!")
	}()

	time.Sleep(100 * time.Millisecond)
}

// Exercise 2: Multiple goroutines
func exercise2() {
	fmt.Println("\n=== Exercise 2: Multiple Goroutines ===")

	// TODO: Start 5 goroutines, each printing its number (1-5)
	// Hint: Use a loop and pass i as parameter

	for i := 1; i <= 5; i++ {
		go func(num int) {
			fmt.Printf("Goroutine %d\n", num)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
}

// Exercise 3: Anonymous goroutine with closure
func exercise3() {
	fmt.Println("\n=== Exercise 3: Closure ===")

	name := "Alice"
	score := 100

	// TODO: Start a goroutine that prints name and score
	// Use closure (access variables from outer scope)

	go func() {
		fmt.Printf("Player: %s, Score: %d\n", name, score)
	}()

	time.Sleep(100 * time.Millisecond)
}

// Exercise 4: WaitGroup
func exercise4() {
	fmt.Println("\n=== Exercise 4: Using WaitGroup ===")

	var wg sync.WaitGroup

	// TODO: Start 3 workers using WaitGroup
	// Each worker should:
	// 1. Be registered with wg.Add(1)
	// 2. Print "Worker X starting"
	// 3. Sleep 100ms
	// 4. Print "Worker X done"
	// 5. Call wg.Done()

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	// TODO: Wait for all workers
	wg.Wait()
	fmt.Println("All workers complete")
}

// Exercise 5: Simulate read/write pumps
func exercise5() {
	fmt.Println("\n=== Exercise 5: Read/Write Pumps ===")

	messages := make(chan string, 5)

	// TODO: Start readPump goroutine
	// Reads from messages channel and prints
	go func() {
		fmt.Println("readPump started")
		for msg := range messages {
			fmt.Printf("readPump received: %s\n", msg)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("readPump ended")
	}()

	// TODO: Start writePump goroutine
	// Sends 3 messages to the channel, then closes it
	go func() {
		fmt.Println("writePump started")
		messages <- "message 1"
		time.Sleep(150 * time.Millisecond)
		messages <- "message 2"
		time.Sleep(150 * time.Millisecond)
		messages <- "message 3"
		close(messages)
		fmt.Println("writePump ended")
	}()

	time.Sleep(600 * time.Millisecond)
}

// Exercise 6: Race condition (see the problem)
func exercise6() {
	fmt.Println("\n=== Exercise 6: Race Condition (The Problem) ===")

	counter := 0

	// Start 100 goroutines, each incrementing counter
	for i := 0; i < 100; i++ {
		go func() {
			counter++ // ❌ NOT SAFE!
		}()
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Counter: %d (should be 100, but might not be!)\n", counter)
	fmt.Println("This is a RACE CONDITION - we'll fix it in Lesson 7 (Mutexes)")
}

// Bonus: Worker pool
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO: Receive jobs from jobs channel
	// TODO: Double each job and send to results
	for job := range jobs {
		results <- job * 2
	}
}

func bonusExercise() {
	fmt.Println("\n=== Bonus: Worker Pool ===")

	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup

	// TODO: Start 3 workers
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// TODO: Send 10 jobs
	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs)

	// Close results when all workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// TODO: Collect results
	sum := 0
	for result := range results {
		sum += result
	}

	fmt.Printf("Sum of results: %d (should be 110)\n", sum)
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   Lesson 6: Goroutine Exercises        ║")
	fmt.Println("╚════════════════════════════════════════╝")

	exercise1()
	exercise2()
	exercise3()
	exercise4()
	exercise5()
	exercise6()
	bonusExercise()

	fmt.Println("\n✅ All exercises complete!")
	fmt.Println("\nNext: Fill out SUMMARY.md, then move to Lesson 4 (Channels)")
}

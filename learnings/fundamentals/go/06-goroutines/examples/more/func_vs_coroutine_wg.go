package main

import (
	"fmt"
	"sync"
	"time"
)

func count(name string, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when function exits

	for i := 1; i <= 5; i++ {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("=== Concurrent with WaitGroup ===")

	wg.Add(1) // We're starting 2 goroutines

	go count("Goroutine-1", &wg)
	wg.Add(1)
	go count("Goroutine-2", &wg)

	wg.Wait() // Block here until both goroutines call Done()

	fmt.Println("\nAll goroutines finished!")
}

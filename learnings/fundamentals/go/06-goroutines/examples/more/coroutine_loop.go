package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== WITHOUT variable capture (BUGGY) ===")

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 1; i <= 5; i++ {
		go func() {
			time.Sleep(10 * time.Millisecond) // Force goroutines to wait
			fmt.Printf("Goroutine-%d\n", i)   // NOW they all read i
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("\n=== WITH variable capture (CORRECT) ===")

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		i := i // Capture the variable
		go func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("Goroutine-%d\n", i)
			wg.Done()
		}()
	}
	wg.Wait()
}

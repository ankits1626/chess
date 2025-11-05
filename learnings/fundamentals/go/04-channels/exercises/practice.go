package main

import (
	"fmt"
	"time"
)

// Exercise 1: Basic Send and Receive
func exercise1() {
	fmt.Println("\n=== Exercise 1: Basic Channel Operations ===")

	// TODO: Create a channel for integers
	numbers := make(chan int)

	// TODO: Send 42 to the channel in a goroutine
	go func() {
		numbers <- 42
	}()

	// TODO: Receive from the channel and print it
	num := <-numbers
	fmt.Printf("Received: %d\n", num)
}

// Exercise 2: Buffered Channel
func exercise2() {
	fmt.Println("\n=== Exercise 2: Buffered Channel ===")

	// TODO: Create a buffered channel with capacity 3
	messages := make(chan string, 3)

	// TODO: Send 3 messages WITHOUT using a goroutine
	messages <- "first"
	messages <- "second"
	messages <- "third"

	// TODO: Receive and print all 3 messages
	fmt.Println(<-messages)
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}

// Exercise 3: Channel Closure Detection
func exercise3() {
	fmt.Println("\n=== Exercise 3: Detect Closed Channel ===")

	messages := make(chan string, 2)

	go func() {
		messages <- "hello"
		messages <- "world"
		close(messages) // Close the channel
	}()

	time.Sleep(100 * time.Millisecond)

	// TODO: Receive messages until channel is closed
	// Use the 'ok' pattern
	for {
		msg, ok := <-messages
		if !ok {
			fmt.Println("Channel closed")
			break
		}
		fmt.Printf("Received: %s\n", msg)
	}
}

// Exercise 4: Range Over Channel
func exercise4() {
	fmt.Println("\n=== Exercise 4: Range Over Channel ===")

	numbers := make(chan int)

	go func() {
		// TODO: Send numbers 1 to 5, then close the channel
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers)
	}()

	// TODO: Use 'range' to receive all numbers
	sum := 0
	for num := range numbers {
		sum += num
	}

	fmt.Printf("Sum: %d (should be 15)\n", sum)
}

// Exercise 5: Worker
func worker(id int, jobs <-chan int, results chan<- int) {
	// TODO: Receive jobs from 'jobs' channel
	// TODO: Double each job and send to 'results' channel
	for job := range jobs {
		results <- job * 2
	}
}

func exercise5() {
	fmt.Println("\n=== Exercise 5: Worker Pattern ===")

	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// TODO: Start 2 workers
	go worker(1, jobs, results)
	go worker(2, jobs, results)

	// TODO: Send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// TODO: Receive 5 results
	for r := 1; r <= 5; r++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

// Exercise 6: Build a Simple Hub
type SimpleHub struct {
	// TODO: Add these channels:
	// - messages: buffered channel (size 10) for []byte
	// - done: unbuffered channel for bool
	messages chan []byte
	done     chan bool
}

func newSimpleHub() *SimpleHub {
	return &SimpleHub{
		messages: make(chan []byte, 10),
		done:     make(chan bool),
	}
}

func (h *SimpleHub) run() {
	fmt.Println("Hub started")
	for {
		select {
		case msg := <-h.messages:
			// TODO: Print the message
			fmt.Printf("Hub received: %s\n", string(msg))

		case <-h.done:
			// TODO: Print "Hub stopping" and return
			fmt.Println("Hub stopping")
			return
		}
	}
}

func exercise6() {
	fmt.Println("\n=== Exercise 6: Simple Hub ===")

	hub := newSimpleHub()
	go hub.run()

	// TODO: Send 3 messages to the hub
	hub.messages <- []byte("Message 1")
	hub.messages <- []byte("Message 2")
	hub.messages <- []byte("Message 3")

	time.Sleep(200 * time.Millisecond)

	// TODO: Signal the hub to stop
	hub.done <- true

	time.Sleep(100 * time.Millisecond)
}

// Bonus Exercise: Ping-Pong
func bonusExercise() {
	fmt.Println("\n=== Bonus: Build Ping-Pong ===")

	// TODO: Create two channels: ping and pong
	ping := make(chan string)
	pong := make(chan string)

	// TODO: Create player 1 goroutine
	// Receives from ping, sends to pong
	go func() {
		for i := 0; i < 3; i++ {
			msg := <-ping
			fmt.Printf("Player 1 got: %s\n", msg)
			pong <- "pong"
		}
	}()

	// TODO: Create player 2 goroutine
	// Receives from pong, sends to ping (except last time)
	go func() {
		for i := 0; i < 3; i++ {
			msg := <-pong
			fmt.Printf("Player 2 got: %s\n", msg)
			if i < 2 {
				ping <- "ping"
			}
		}
	}()

	// TODO: Start the game by sending first "ping"
	ping <- "ping"

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Game over!")
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   Lesson 4: Channel Exercises          ║")
	fmt.Println("╚════════════════════════════════════════╝")

	exercise1()
	exercise2()
	exercise3()
	exercise4()
	exercise5()
	exercise6()
	bonusExercise()

	fmt.Println("\n✅ All exercises complete!")
	fmt.Println("\nNext: Fill out SUMMARY.md with your learnings")
}

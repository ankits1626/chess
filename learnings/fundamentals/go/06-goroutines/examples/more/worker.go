package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)
		time.Sleep(1 * time.Second)
		results <- job * 2
		fmt.Printf("From Worker %d: Received job %d result %d\n", id, job, job*2)
	}
}

func main() {
	numJobs := 20
	numWorkers := 2
	jobs := make(chan int)
	results := make(chan int)

	for w := 0; w < numWorkers; w++ {
		go worker(w, jobs, results)
	}

	go func() {
		for i := 0; i < numJobs; i++ {
			fmt.Printf("Pushing job %d\n", i)
			jobs <- i
		}
		close(jobs)
	}()

	for i := 0; i < numJobs; i++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func doWork(id int, ch chan string, wg *sync.WaitGroup) {
	// defer wg.Done()
	ch <- fmt.Sprintf("Work %d done", id)
}

func main() {
	ch := make(chan string, 3)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		// wg.Add(1)
		go doWork(i, ch, &wg)
	}

	// go func() {
	// 	wg.Wait()
	// 	close(ch)
	// }()
	// for msg := range ch {
	// 	fmt.Println(msg)
	// }
	// m1 := <-ch
	// fmt.Println("Messsage 1: ", m1)

	// m2 := <-ch
	// fmt.Println("Messsage 2: ", m2)

	// m3 := <-ch
	// fmt.Println("Messsage 3: ", m3)
	time.Sleep(1 * time.Second)
	fmt.Println("Main complete")
}

package main

import (
	"fmt"
	"sync"
)

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Task %d is starting\n", id)
	// Simulate some work
	for range 5 {
		fmt.Printf("Task %d is working...\n", id)
	}
	fmt.Printf("Task %d is done\n", id)
}

func main() {
	fmt.Println("This is the main package for the goroutines example.")

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go task(i, &wg)
	}

	// Wait for the goroutine to finish
	fmt.Println("Waiting for the goroutine to finish...")

	wg.Wait()
	fmt.Println("Goroutine has finished execution")
}

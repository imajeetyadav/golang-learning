package main

import "fmt"

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println("Iteration:", i)
	}
	for j := 5; j > 0; j-- {
		fmt.Println("Countdown:", j)
	}
	numbers := []int{1, 2, 3, 4, 5}
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}
	// Using a for loop with a condition
	i := 0
	for i < 5 {
		fmt.Println("Condition Loop:", i)
		i++
	}
}

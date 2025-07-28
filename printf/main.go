package main

import "fmt"

func main() {
	// Create a map with string keys and int values
	myMap := map[string]int{
		"apple":  5,
		"banana": 3,
		"orange": 8,
	}

	// Iterate over the map using a for range loop
	for key, value := range myMap {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}
}

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(greet("World"))
	fmt.Println(time.Now().Year())
}

func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
